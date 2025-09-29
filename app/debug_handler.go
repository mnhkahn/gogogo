package app

import (
	"bytes"
	"html/template"
	"os"
	"runtime"
	"sort"

	"github.com/mnhkahn/gogogo/logger"
	"github.com/mnhkahn/gogogo/util"
)

// GoAppHandler is a handler to show debug information.
func GoAppHandler(c *Context) error {
	res := make(map[string]interface{}, 4)
	res["app.VERSION"] = VERSION
	res["app.BUILD"] = BUILD
	res["app.BRANCH"] = BRANCH
	res["go.version"] = runtime.Version()
	res["os.args"] = os.Args
	res["os.env"] = os.Environ()
	res["os.goroutine"] = util.Goroutine()
	res["os.pwd"], _ = os.Getwd()
	res["log.level"] = logger.StdLogger.GetLevel()
	c.JSON(res)

	return nil
}

var debugRouterTpl = `
<html lang="en">
<head>
   <title>Router</title>
</head>
<body>
<h3>Router</h3>
<ul>
{{range $i, $r := .Routers}}
    <li><a href="{{$r}}">{{$r}}</a></li>
{{end}}
</ul>
</body>
`

// LogLevelHandler is a handler to set log level for StdLogger.
func LogLevelHandler(c *Context) error {
	l, err := c.GetInt("level")
	if err != nil {
		return err
	}
	res := logger.StdLogger.SetLevel(l)
	c.JSON(res)

	return nil
}

func StatHandler(c *Context) error {
	var buf bytes.Buffer
	tpl := template.New("statTpl")
	tpl = template.Must(tpl.Parse(statTpl))
	stats := make([]*Stat, 0, len(DefaultHandler.Stats))
	for _, v := range DefaultHandler.Stats {
		stats = append(stats, v)
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Cnt == stats[j].Cnt {
			return stats[i].Url > stats[j].Url
		}
		return stats[i].Cnt > stats[j].Cnt
	})
	err := tpl.ExecuteTemplate(&buf, "statTpl", struct {
		Stats []*Stat
	}{
		stats,
	})
	if err != nil {
		return err
	}

	c.WriteBytes(buf.Bytes())

	return nil
}

var statTpl = `
<html lang="en">
    <head>
        <title>Statistics</title>
    </head>
    <body>
        <h3>Statistics</h3>
        <script>
            // 定义静态资源的扩展名
            const staticExtensions = [".css", ".js", ".ico", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".woff", ".woff2", ".ttf", ".eot", ".env", ".env.production", "phpinfo", ".django"];

            function toggleStaticResources() {
                const checkbox = document.getElementById("filterStatic");
                const rows = document.querySelectorAll("table tr:not(:first-child)"); // 排除表头

                rows.forEach((row) => {
                    const urlCell = row.querySelector("td:first-child");
                    const url = urlCell.textContent.trim();

                    // 检查URL是否包含静态资源扩展名
                    const isStatic = staticExtensions.some((ext) => url.endsWith(ext));

                    // 根据复选框状态和是否为静态资源来显示或隐藏行
                    if (checkbox.checked && isStatic) {
                        row.style.display = "none";
                    } else {
                        row.style.display = "";
                    }
                });
            }
			function toggleSuccessRequests() {
                const checkbox = document.getElementById("filterSuccess");
                const rows = document.querySelectorAll("table tr:not(:first-child)"); // 排除表头

                rows.forEach((row) => {
                    const urlCell = row.querySelector("td:nth-child(3)");
                    const url = urlCell.textContent.trim();

                    // 检查URL是否包含静态资源扩展名
                    const isError = parseInt(url, 10) >= 300;

                    // 根据复选框状态和是否为静态资源来显示或隐藏行
                    if (checkbox.checked && isError) {
                        row.style.display = "none";
                    } else {
                        row.style.display = "";
                    }
                });
            }
        </script>
        <div>
            <input type="checkbox" id="filterStatic" onchange="toggleStaticResources()" />
            <label for="filterStatic">Exclude static resources</label>
			<input type="checkbox" id="filterSuccess" onchange="toggleSuccessRequests()" />
            <label for="filterSuccess">Exclude error requests</label>
        </div>
        <table style="width: 100%">
            <tr>
                <th>Url Path</th>
                <th>Count</th>
                <th>Status Code</th>
                <th>Sum Elapse</th>
                <th>AvgTime Elapse</th>
            </tr>
            {{range $stat := .Stats}}
            <tr>
                <td align="center" width="60%">{{$stat.Url}}</td>
                <td align="center">{{$stat.Cnt}}</td>
                <td align="center">{{$stat.StatusCode}}</td>
                <td align="center">{{$stat.SumTime}}</td>
                <td align="center">{{$stat.AvgTime}}</td>
            </tr>
            {{end}}
        </table>
    </body>
</html>
`
