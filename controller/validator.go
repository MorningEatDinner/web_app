package controller

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/thedevsaddam/govalidator"
)

// 定义一个全局翻译器T
var trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) { //local就是你想要翻译成为什么
	// locale 通常取决于 http 请求头的 'Accept-Language'

	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取json tag的自定义方法
		// 这个部分的目的就是说将返回的信息读取为json格式希望得到的那个字段， 而不是程序结构体中标明的那个字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

//type SignUpParam struct {
//	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
//	Name       string `json:"name" binding:"required"`
//	Email      string `json:"email" binding:"required,email"`
//	Password   string `json:"password" binding:"required"`
//	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
//}

//func main() {
//	if err := InitTrans("zh"); err != nil {
//		fmt.Printf("init trans failed, err:%v\n", err)
//		return
//	}
//
//	r := gin.Default()
//
//	r.POST("/signup", func(c *gin.Context) {
//		var u SignUpParam
//		if err := c.ShouldBind(&u); err != nil {
//			// 获取validator.ValidationErrors类型的errors
//			errs, ok := err.(validator.ValidationErrors)
//			if !ok {
//				// 非validator.ValidationErrors类型错误直接返回
//				c.JSON(http.StatusOK, gin.H{
//					"msg": err.Error(),
//				})
//				return
//			}
//			// validator.ValidationErrors类型错误则进行翻译
//			c.JSON(http.StatusOK, gin.H{
//				"msg": errs.Translate(trans),
//			})
//			return
//		}
//		// 保存入库等具体业务逻辑代码...
//
//		c.JSON(http.StatusOK, "success")
//	})
//
//	_ = r.Run(":8999")
//}

// 去除提示信息中的结构体名称
/*
   "msg": {
       "ParamSignUp.re_password": "re_password为必填字段"
   }
*/
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

type ValidatorFunc func(interface{}, *gin.Context) map[string][]string

func Validate(ctx *gin.Context, obj interface{}, handler ValidatorFunc) bool {
	if err := ctx.ShouldBind(obj); err != nil {
		// 参数解析失败
		//ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		//	"error": err.Error(),
		//})
		BadRequest(ctx, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		// 打印错误信息
		fmt.Println(err.Error())
		// 错误之后中断请求
		return false
	}
	// 绑定验证器
	errs := handler(obj, ctx)
	if len(errs) > 0 {
		// 验证失败
		//ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		//	"errors": errs,
		//})
		ValidationError(ctx, errs) // 如果参数解析失败， 直接在这里返回响应
		return false
	}

	return true
}

// ValidateSignupPhoneExist：验证器函数， 验证数据是否符合要求
func ValidateSignupPhoneExist(data interface{}, ctx *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"phone": []string{"required", "digits:11"}, // 定义每个字段需要满足的规则是什么
	}
	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项，参数名称 phone", // 如果不满足这个字段的要求， 那么会返回这个信息
			"digits:手机号长度必须为 11 位的数字",
		},
	}

	return validate(data, rules, messages)
}

// validate：根据规则对于传输进来的数据进行验证
func validate(data interface{}, rules, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Data:          data,     // 请求验证的结构
		Rules:         rules,    // 加入这个tag需要满足什么功能
		Messages:      messages, // 如果错误需要返回的信息是什么
		TagIdentifier: "valid",  // 在结构体中的tag是什么
	}
	return govalidator.New(opts).ValidateStruct()
}
