package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

func createHTMLFile(section string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Проверка существования HTML файла
	_, err := os.Stat(fmt.Sprintf("./src/html/partials/sections/%s.html", section))
	if err == nil {
		fmt.Println(fmt.Errorf("HTML file for section %s already exists", section))
		return
	}

	// Создание HTML файла в директории ./src/html/partials/sections
	htmlFile, err := os.Create(fmt.Sprintf("./src/html/partials/sections/%s.html", section))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer htmlFile.Close()

	// Запись в HTML файл
	htmlContent := fmt.Sprintf("{%% set section = %s %%}\n\n<section class=\"{{section.id}} section\">\n  <div class=\"container\">\n   <h2 class=\"section-title\">{{section.title}}</h2>\n  </div>\n</section>", section)
	_, err = htmlFile.WriteString(htmlContent)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func createSCSSFile(section string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Проверка существования SCSS файла
	_, err := os.Stat(fmt.Sprintf("./src/assets/scss/sections/_%s.scss", section))
	if err == nil {
		fmt.Println(fmt.Errorf("SCSS file for section %s already exists", section))
		return
	}

	// Создание SCSS файла в директории ./src/assets/scss/sections
	scssFile, err := os.Create(fmt.Sprintf("./src/assets/scss/sections/_%s.scss", section))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer scssFile.Close()

	// Запись в SCSS файл
	scssContent := fmt.Sprintf(".%s { }", section)
	_, err = scssFile.WriteString(scssContent)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	// Чтение JSON файла
	data, err := os.ReadFile("gulp-starter.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Декодирование JSON в массив строк
	var sections []string
	err = json.Unmarshal(data, &sections)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	for _, section := range sections {
		wg.Add(2)
		go createHTMLFile(section, &wg)
		go createSCSSFile(section, &wg)
	}

	wg.Wait()

	// Чтение HTML файла
	indexHTML, err := os.ReadFile("./src/html/pages/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Создание строки для вставки в HTML
	var includes []string
	for _, section := range sections {
		includes = append(includes, fmt.Sprintf("{%% include 'sections/%s.html' %%}", section))
	}
	insert := strings.Join(includes, "\n")

	// Вставка строки в HTML
	newHTML := strings.Replace(string(indexHTML), "<!--  -->", insert, 1)

	// Запись обновленного HTML обратно в файл
	err = os.WriteFile("./src/html/pages/index.html", []byte(newHTML), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Чтение SCSS файла
	mainSCSS, err := os.ReadFile("./src/assets/scss/main.scss")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Создание строки для вставки в SCSS
	var imports []string
	for _, section := range sections {
		importLine := fmt.Sprintf("@import 'sections/%s';", section)
		imports = append(imports, importLine)

	}

	insertSCSS := strings.Join(imports, "\n")

	// Вставка строки в SCSS
	newSCSS := strings.Replace(string(mainSCSS), "/* <!--  --> */", insertSCSS, 1)

	// Запись обновленного SCSS обратно в файл
	err = os.WriteFile("./src/assets/scss/main.scss", []byte(newSCSS), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Files created successfully")
}
