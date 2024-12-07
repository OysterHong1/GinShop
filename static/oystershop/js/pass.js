(function($){
    $(function(){
        loginApp.init();
    })
    var loginApp={
        init:function(){
            this.getCaptcha();
            this.captchaImgChage();
            this.initRegisterStep1();
			this.initRegisterStep2();
			this.initRegisterStep3();
			this.initDoLogin();
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
        },initRegisterStep1:function(){
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
        },initRegisterStep2:function(){
			$(function () {
				var timer = 10;
				function Countdown() {
					if (timer >= 1) {
						timer -= 1;
						$("#sendCode").attr('disabled', true);
						$("#sendCode").html('重新发送(' + timer + ')');
						setTimeout(function () {
							Countdown();
						}, 1000);
					} else {
						$("#sendCode").attr('disabled', false)
						$("#sendCode").html('重新发送');
					}
				}
				Countdown();

				$("#sendCode").click(function () {
					timer = 10;
					Countdown();
					var phone = $("#phone").val()
					var verifyCode = $("#verifyCode").val()
					var captchaId = "resend"

					//重新请求接口发送短信
					$.get("/pass/sendCode", { "phone": phone, "verifyCode": verifyCode, "captchaId": captchaId }, function (response) {
						console.log(response)
					})
				})
			})

			//验证验证码		
			$(function () {

				$("#nextStep").click(function (e) {
					$(".error").html()
					var sign = $('#sign').val();
					var smsCode = $('#smsCode').val();

					$.get('/pass/validateSmsCode', { sign: sign, smsCode: smsCode }, function (response) {
						console.log(response)
						if (response.success == true) {
							location.href = "/pass/registerStep3?sign=" + sign + "&smsCode=" + smsCode
						} else {
							$(".error").html("Error：" + response.message)
						}
					})

				})
				
				$("#returnButton").click(function(){
						location.href="/pass/registerStep1"
				})

			})
		},initRegisterStep3:function(){
			$(function(){
				$("#form").submit(function(){
					$(".error").html("")
					var password=$('#password').val();
					var rpassword=$('#rpassword').val();

					if(password.length<6){
						$(".error").html("Error：密码的长度不能小于6位")						
						return false;
					}					
					if(password!=rpassword){							
						$(".error").html("Error：密码和确认密码不一致")
						return false;
					}
					return true;

				})
			})
		},initDoLogin:function(){
			var _that=this;
			$("#doLogin").click(function(e){
				$(".error").html("")
				var phone=$('#phone').val();
				var password= $('#password').val();			
				var captchaId = $('#captchaId').val();
				var captchaVal = $("#captchaVal").val();
				var prevPage = $("#prevPage").val();
				var reg =/^[\d]{11}$/;
				if(!reg.test(phone)){
					$(".error").html('Error:手机号输入错误');					
					return false;
				}
				if(password.length<6){
					$(".error").html('Error:密码长度不合法');						
					return false;
				}

				if(captchaVal.length<4){						
					$(".error").html('Error:验证码长度不合法');	
					return false;
				}
				//ajax请求	 														
				$.post('/pass/doLogin',{phone:phone,password:password,captchaVal:captchaVal,captchaId:captchaId},function(response){							
					console.log(response);
					if(response.success==true){						
						if(prevPage == ""){
							location.href="/";
						}else{
							location.href=prevPage;
						}
					}else{															
						$(".error").html("Error：" + response.message + ",请重新输入!")
						//改变验证码	
						_that.getCaptcha();				
					}					
				})
			})
		}
    }
})($)

