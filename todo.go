package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Date       string
	Time_start string
	Time_end   string
	To_do      string
	Which_do   string
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func CreateTasks() (string, string, string, int, int, int) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	var pro_sum_toushi int
	var pro_sum_shouhi int
	var pro_sum_rouhi int

	product := []Product{}

	/*データベースから「投資」の項目を取得し、時間の合計値を取る*/
	db.Where("which_do = ?", "投資").Find(&product)
	//	db.Where("which_do = ?", "投資").Delete(&product)
	for _, pro := range product {

		time_start_head := pro.Time_start[:2]
		time_start_bottom := pro.Time_start[3:]
		time_end_head := pro.Time_end[:2]
		time_end_bottom := pro.Time_end[3:]

		time_start_head_int, _ := strconv.Atoi(time_start_head)
		time_start_bottom_int, _ := strconv.Atoi(time_start_bottom)
		time_end_head_int, _ := strconv.Atoi(time_end_head)
		time_end_bottom_int, _ := strconv.Atoi(time_end_bottom)

		time_start := time_start_head_int*60 + time_start_bottom_int
		time_end := time_end_head_int*60 + time_end_bottom_int

		if time_start > time_end {
			time_end = time_end + 24*60
		}

		time_differ := time_end - time_start

		//		fmt.Println(time_differ)

		pro_sum_toushi += time_differ

	}

	/*データベースから「消費」の項目を取得し、時間の合計値を取る*/
	db.Where("which_do = ?", "消費").Find(&product)
	//	db.Where("which_do = ?", "投資").Delete(&product)
	for _, pro := range product {
		//pro2 = append(pro2, pro.Time_start)
		//pro2 = append(pro2, pro.Time_end)

		time_start_head := pro.Time_start[:2]
		time_start_bottom := pro.Time_start[3:]
		time_end_head := pro.Time_end[:2]
		time_end_bottom := pro.Time_end[3:]

		time_start_head_int, _ := strconv.Atoi(time_start_head)
		time_start_bottom_int, _ := strconv.Atoi(time_start_bottom)
		time_end_head_int, _ := strconv.Atoi(time_end_head)
		time_end_bottom_int, _ := strconv.Atoi(time_end_bottom)

		time_start := time_start_head_int*60 + time_start_bottom_int
		time_end := time_end_head_int*60 + time_end_bottom_int

		if time_start > time_end {
			time_end = time_end + 24*60
		}

		time_differ := time_end - time_start

		//		fmt.Println(time_differ)

		pro_sum_shouhi += time_differ

	}
	/*データベースから「浪費」の項目を取得し、時間の合計値を取る*/
	db.Where("which_do = ?", "浪費").Find(&product)
	//	db.Where("which_do = ?", "投資").Delete(&product)
	for _, pro := range product {

		time_start_head := pro.Time_start[:2]
		time_start_bottom := pro.Time_start[3:]
		time_end_head := pro.Time_end[:2]
		time_end_bottom := pro.Time_end[3:]

		time_start_head_int, _ := strconv.Atoi(time_start_head)
		time_start_bottom_int, _ := strconv.Atoi(time_start_bottom)
		time_end_head_int, _ := strconv.Atoi(time_end_head)
		time_end_bottom_int, _ := strconv.Atoi(time_end_bottom)

		time_start := time_start_head_int*60 + time_start_bottom_int
		time_end := time_end_head_int*60 + time_end_bottom_int

		if time_start > time_end {
			time_end = time_end + 24*60
		}

		time_differ := time_end - time_start

		//		fmt.Println(time_differ)

		pro_sum_rouhi += time_differ

	}

	/*順位付けを行う構造体*/
	type Rank struct {
		What string
		Sum  int
		Moji string
	}

	rank := []Rank{
		{What: "pro_sum_toushi", Sum: pro_sum_toushi, Moji: "投資"},
		{What: "pro_sum_shouhi", Sum: pro_sum_shouhi, Moji: "消費"},
		{What: "pro_sum_rouhi", Sum: pro_sum_rouhi, Moji: "浪費"},
	}

	/*時間順に並べ替えを行う*/
	sort.Slice(rank, func(i, j int) bool { return rank[i].Sum < rank[j].Sum })
	fmt.Printf("1位:%+v\n", rank[2].Moji)
	fmt.Printf("2位:%+v\n", rank[1].Moji)
	fmt.Printf("3位:%+v\n", rank[0].Moji)
	fmt.Printf("1位:%+v\n", rank[2].Sum)
	fmt.Printf("2位:%+v\n", rank[1].Sum)
	fmt.Printf("3位:%+v\n", rank[0].Sum)

	return rank[2].Moji, rank[1].Moji, rank[0].Moji, rank[2].Sum, rank[1].Sum, rank[0].Sum

}

/*フロントエンドへ値を渡す構造体*/
type Data struct {
	Rank1_moji string
	Rank2_moji string
	Rank3_moji string
	Rank1_time int
	Rank2_time int
	Rank3_time int
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}

	e := echo.New()

	e.Renderer = t

	e.File("/", "public/index.html")
	e.GET("/writeout", func(c echo.Context) error {
		/*cはechoの変数で使っているので使えない*/
		a, b, d, e, f, g := CreateTasks()
		/*CreateTasksの返り値(種別と時間)が順位順になってa,b,cd,e,f,gに代入*/
		var data Data
		data.Rank1_moji = a
		data.Rank2_moji = b
		data.Rank3_moji = d
		data.Rank1_time = e
		data.Rank2_time = f
		data.Rank3_time = g

		return c.Render(http.StatusOK, "index2.html", data)
	})

	e.Start(":8080")

}
