(function($){
    var app={
        init:function(){             
            this.changeCartNum();       
            this.deleteConfirm();
            this.initCheckBox();
            this.isCheckedAll();
            this.initChekOut();
        },
        deleteConfirm:function(){
            $('.delete').click(function(){    
                var flag=confirm('您确定要删除吗?');    
                return flag;    
            })    
        },     
        initChekOut(){
            $(function(){	
                    $("#checkout").click(function(){	
                        var allPrice=parseFloat($("#allPrice").html());	
                        if(allPrice==0){
                            alert('购物车没有选中去结算的商品')
                        }else{
                            location.href="/buy/checkout";
                        }	
                    })
            })
        },
        initCheckBox(){
            //全选按钮点击
            $("#checkAll").click(function() {               
                if (this.checked) {
                    $(":checkbox").prop("checked", true);
                    //让cookie中商品的checked属性都等于true                
                    $.get('/cart/changeAllCart?flag=1',function(response){                 
                        if(response.success){
                            $("#allPrice").html(response.allPrice+"元") 
                        }
                    })
                }else {
                    $(":checkbox").prop("checked", false);      
                     //让cookie中商品的checked属性都等于false
                    $.get('/cart/changeAllCart?flag=0',function(response){                 
                        if(response.success){
                            $("#allPrice").html(response.allPrice+"元") 
                        }
                    })                           
                }               
            });    

            //点击单个选择框按钮的时候触发       
            var _that=this;            
            $(".cart_list :checkbox").click(function() {                         
                _that.isCheckedAll();

                var goods_id=$(this).attr("goods_id")
                var goods_color=$(this).attr("goods_color")
                $.get('/cart/changeOneCart?goods_id='+goods_id+'&goods_color='+goods_color,function(response){                 
                    if(response.success){
                        $("#allPrice").html(response.allPrice+"元") 
                    }
                })
               

            });   //注意：this指向
        },
        //判断全选是否选择
        isCheckedAll(){             
            var allNum = $(".cart_list :checkbox").size();//checkbox总个数
            var checkedNum = 0;  
            $(".cart_list :checkbox").each(function () {  
                if($(this).prop("checked")==true){
                    checkedNum++;
                }
            });
            console.log(allNum,checkedNum)
            if(allNum==checkedNum){//全选
                $("#checkAll").prop("checked",true);
            }else{//不全选
                $("#checkAll").prop("checked",false);
            }
        }, 
        changeCartNum(){
            $('.decCart').click(function(){
                var goods_id=$(this).attr("goods_id")
                var goods_color=$(this).attr("goods_color")
                var _that=this;
                $.get('/cart/decCart?goods_id='+goods_id+'&goods_color='+goods_color,function(response){
                    console.log(response)
                    if(response.success){
                        $("#allPrice").html(response.allPrice+"元")
                        $(_that).siblings(".input_center").find("input").val(response.num)
                        $(_that).parent().parent().siblings(".totalPrice").html(response.currentPrice+"元")
                        
                    }
                })

            });

            $('.incCart').click(function(){
                var goods_id=$(this).attr("goods_id")
                var goods_color=$(this).attr("goods_color")
                var _that=this;
                

                $.get('/cart/incCart?goods_id='+goods_id+'&goods_color='+goods_color,function(response){
                    console.log(response)
                    if(response.success){
                        $("#allPrice").html(response.allPrice+"元")
                        $(_that).siblings(".input_center").find("input").val(response.num)
                        $(_that).parent().parent().siblings(".totalPrice").html(response.currentPrice+"元")
                        
                    }
                })

            });
        }
    }

    $(function(){
        app.init();
    })    
})($)
