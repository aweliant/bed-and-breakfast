package render

import (
	"bytes"
	"fmt"
	"github.com/aweliant/bed-and-breakfast/pkg/config"
	"github.com/aweliant/bed-and-breakfast/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var cfg *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	cfg = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

//renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	//w 是用以存return的webpage的
	tc := map[string]*template.Template{}
	if cfg.UseCache {
		// get the template cache from the app config
		tc = cfg.TemplateCache
	} else {
		tc, _ = GenerateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("cpuld not get tempalte from temoklate cacghe")
	}
	buf := new(bytes.Buffer)
	//Why do we need this?用缓存来提升IO性能。缓冲区在创建时就被分配内存，这块内存区域一直被重用，可以减少动态分配和回收内存的次数
	// 缓冲byte类型的缓冲器.具有读写方法和可变大小的字节存储功能。
	//和var buf *bytes.Buffer一样？
	//到底为何用buffer,Trevor给的答案是：
	//https://www.udemy.com/course/building-modern-web-applications-with-go/learn/lecture/22870657#questions/14460376
	//I write to to the buffer, instead of straight to the ResponseWriter because I can check to see if there's an error,
	//and determine where it comes from more easily.
	td = AddDefaultData(td)
	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	//parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	//err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return
	}
}

func GenerateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page) //.Funcs must be called before the template is parsed.
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		//之所以要这么检查一下len(matches) > 0，是因为ParseGlob要求必须能匹配上至少一个。
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl") //批量解析。这里也可以用ParseFiles以具体模板路径为参数。
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
