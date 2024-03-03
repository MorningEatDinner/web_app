package controller

import (
	"fmt"
	"github.com/xiaorui/web_app/dao/redis"
	"github.com/xiaorui/web_app/models"
	"github.com/xiaorui/web_app/pkg/captcha"
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

// ValidateSignupEmailExist：验证器函数， 验证数据是否符合要求
func ValidateSignupEmailExist(data interface{}, ctx *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"email": []string{"required", "min:4", "max:30", "email"}, // 定义每个字段需要满足的规则是什么
	}
	messages := govalidator.MapData{
		"email": []string{
			"required:Email 为必填项",
			"min:Email 长度需大于 4",
			"max:Email 长度需小于 30",
			"email:Email 格式不正确，请提供有效的邮箱地址",
		},
	}

	return validate(data, rules, messages)
}

func ValidatePhoneCodeRequest(data interface{}, ctx *gin.Context) map[string][]string {
	// 1. 规则
	rules := govalidator.MapData{
		"phone":          []string{"required", "digits:11"},
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
	}
	// 2. 错误信息
	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项，参数名称 phone",
			"digits:手机号长度必须为 11 位的数字",
		},
		"captcha_id": []string{
			"required:图片验证码的 ID 为必填",
		},
		"captcha_answer": []string{
			"required:图片验证码答案必填",
			"digits:图片验证码长度必须为 6 位的数字",
		},
	}
	// 3. 验证数据
	errs := validate(data, rules, messages)
	// 验证图片验证码有没有问题
	_data := data.(*models.ParamPhoneCode)
	errs = ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs
}

func ValidateEmailCodeRequest(data interface{}, ctx *gin.Context) map[string][]string {
	// 1. 规则
	rules := govalidator.MapData{
		"email":          []string{"required", "min:4", "max:30", "email"},
		"captcha_id":     []string{"required"},
		"captcha_answer": []string{"required", "digits:6"},
	}
	// 2. 错误信息
	messages := govalidator.MapData{
		"email": []string{
			"required:Email 为必填项",
			"min:Email 长度需大于 4",
			"max:Email 长度需小于 30",
			"email:Email 格式不正确，请提供有效的邮箱地址",
		},
		"captcha_id": []string{
			"required:图片验证码的 ID 为必填",
		},
		"captcha_answer": []string{
			"required:图片验证码答案必填",
			"digits:图片验证码长度必须为 6 位的数字",
		},
	}
	// 3. 验证数据
	errs := validate(data, rules, messages)

	// 验证图形验证码是否正确
	_data := data.(*models.ParamEmailCode)
	errs = ValidateCaptcha(_data.CaptchaID, _data.CaptchaAnswer, errs)

	return errs
}

// ValidateSignupUsingPhone: 验证使用手机号码进行注册的参数是否正确
func ValidateSignupUsingPhone(data interface{}, ctx *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"phone":            []string{"required", "digits:11"},
		"name":             []string{"required", "alpha_num", "between:3,20"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
		"verify_code":      []string{"required", "digits:6"},
	}

	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项，参数名称 phone",
			"digits:手机号长度必须为 11 位的数字",
		},
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:用户名格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
		},
		"password": []string{
			"required:密码为必填项",
			"min:密码长度需大于 6",
		},
		"password_confirm": []string{
			"required:确认密码框为必填项",
		},
		"verify_code": []string{
			"required:验证码答案必填",
			"digits:验证码长度必须为 6 位的数字",
		},
	}

	errs := validate(data, rules, messages)
	_data := data.(*models.ParamSignupUsingPhone)
	errs = ValidatePasswordConfirm(_data.Password, _data.PasswordConfirm, errs)
	errs = ValidateKeyCode(_data.Phone, _data.Code, errs)

	return errs
}

// ValidateSignupUsingEmail：进行使用邮箱进行注册的参数验证功能
func ValidateSignupUsingEmail(data interface{}, ctx *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"email":            []string{"required", "min:4", "max:30", "email"},
		"name":             []string{"required", "alpha_num", "between:3,20"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
		"verify_code":      []string{"required", "digits:6"},
	}

	messages := govalidator.MapData{
		"email": []string{
			"required:Email 为必填项",
			"min:Email 长度需大于 4",
			"max:Email 长度需小于 30",
			"email:Email 格式不正确，请提供有效的邮箱地址",
		},
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:用户名格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
		},
		"password": []string{
			"required:密码为必填项",
			"min:密码长度需大于 6",
		},
		"password_confirm": []string{
			"required:确认密码框为必填项",
		},
		"verify_code": []string{
			"required:验证码答案必填",
			"digits:验证码长度必须为 6 位的数字",
		},
	}

	errs := validate(data, rules, messages) // 根据规则来验证参数

	_data := data.(*models.ParamSignUpUsingEmail)
	// 验证验证码是否正确
	errs = ValidateKeyCode(_data.Email, _data.Code, errs)
	// 验证密码是否相同
	errs = ValidatePasswordConfirm(_data.Password, _data.PasswordConfirm, errs)

	return errs
}

// ValidateLoginUsingPhone: 使用手机号码+验证码进行登陆
func ValidateLoginUsingPhoneWithCode(data interface{}, ctx *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"phone":       []string{"required", "digits:11"},
		"verify_code": []string{"required", "digits:6"},
	}
	messages := govalidator.MapData{
		"phone": []string{
			"required:手机号为必填项，参数名称 phone",
			"digits:手机号长度必须为 11 位的数字",
		},
		"verify_code": []string{
			"required:验证码答案必填",
			"digits:验证码长度必须为 6 位的数字",
		},
	}

	errs := validate(data, rules, messages)
	// 进行强制转换之后验证验证码是否正确
	_data := data.(*models.ParamLoginUsingPhoneWithCode)
	errs = ValidateKeyCode(_data.Phone, _data.Code, errs)

	return errs
}

func ValidateCaptcha(captchaID, captchaAnswer string, errs map[string][]string) map[string][]string {
	if ok := captcha.NewCaptcha().VerifyCaptcha(captchaID, captchaAnswer); !ok {
		errs["captcha_answer"] = append(errs["captcha_answer"], "图片验证码错误")
	}
	return errs
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

// ValidatePasswordConfirm: 验证两次输入的手机号码是否相等
func ValidatePasswordConfirm(password, passwordConfirm string, errs map[string][]string) map[string][]string {
	if password != passwordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"], "两次输入的密码不一致")
	}
	return errs
}

// ValidateKeyCode: 验证手机或者邮箱验证码是否正确
func ValidateKeyCode(key, code string, errs map[string][]string) map[string][]string {
	if ok := redis.CheckVerifyCode(key, code); !ok {
		errs["verify_code"] = append(errs["verify_code"], "验证码错误")
	}
	return errs
}
