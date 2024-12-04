(function($){
    $(function(){
        loginApp.init();
    })
    var loginApp={
        init:function(){
            this.getCaptcha()
            this.captchaImgChage()
            this.initRegisterStep1()
        },
        getCaptcha:function(){
            $.get("/pass/captcha?t="+Math.random(),function(response){              
                $("#captchaId").val(response.captchaId)
                $("#captchaImg").attr("src",response.captchaImage)
            })
        },
        captchaImgChage:function(){
            var _that=this;
            $("#captchaImg").click(function(){
                _that.getCaptcha()
            })
        },
        initRegisterStep1:function(){
            var _that=this;
            //发送验证码
			$("#registerButton").click(function () {
				//验证验证码是否正确
				var phone = $('#phone').val();
				var verifyCode = $('#verifyCode').val();
				var captchaId = $("#captchaId").val();
				$(".error").html("")

				var reg = /^[\d]{11}$/;
				if (!reg.test(phone)) {
					$(".error").html("Error：手机号输入错误");
					return false;
				}
				if (verifyCode.length < 2) {
					$(".error").html("Error：图形验证码长度不合法")
					return false;
				}						
				
				$.get("/pass/sendCode",{"phone":phone,"verifyCode":verifyCode,"captchaId":captchaId},function(response){
					console.log(response)
					if (response.success == true) {						
						//跳转到下页面
						location.href="/pass/registerStep2?sign="+response.sign+"&verifyCode="+verifyCode;				
					} else {
						//改变验证码											
						$(".error").html("Error：" + response.message + ",请重新输入!")
						//改变验证码
                        _that.getCaptcha()
												
					}
				})			

			})
        }
    }
})($)

