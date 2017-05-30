define([],function(){
	//字节到gb
	window.bytesToGB=function(bytes){
		var gb=(bytes/1024/1024/1024).toFixed(2).toString();
		return gb ;
	}
	//格式化
	String.prototype.format = function(args) {
	    var result = this;
	    if (arguments.length > 0) {    
	        if (arguments.length == 1 && typeof (args) == "object") {
	            for (var key in args) {
	                if(args[key]!=undefined){
	                    var reg = new RegExp("({" + key + "})", "g");
	                    result = result.replace(reg, args[key]);
	                }
	            }
	        }
	        else {
	            for (var i = 0; i < arguments.length; i++) {
	                if (arguments[i] != undefined) {
	                    //var reg = new RegExp("({[" + i + "]})", "g");//这个在索引大于9时会有问题，谢谢何以笙箫的指出
	　　　　　　　　　　　　var reg= new RegExp("({)" + i + "(})", "g");
	                    result = result.replace(reg, arguments[i]);
	                }
	            }
	        }
	    }
	    return result;
	}
	
	window.AsyncUploadFile=function(option){
		var fd = new FormData(),xhr = new XMLHttpRequest();
			var input = document.createElement('input');
			input.setAttribute('id', 'uploadFile');
			input.setAttribute('type', 'file');
			input.setAttribute('name', 'file');
			document.body.appendChild(input);
			input.style.display = 'none';
			input.click();
			input.onchange = function(){
			if(!option||!option.url){
				option.onFinish&&option.onFinish(false,"url not set!")
				return
			}
			if(input.files[0].size > 5 * 1024 * 1024){
				option.onFinish&&option.onFinish(false,"filesize is too large!")
				return ;
			}
			fd.append('file', input.files[0]);
			//fd.append('filetype', "key");
			if(option.data){
				for(var key in option.data){
					fd.append(key,option.data[key]);
				}
			}
    		xhr.open('post','/uploadfile');
			if(option.onFinish){
				xhr.onreadystatechange = function(){
	           if(xhr.status == 200){
	               if(xhr.readyState == 4){
	                  option.onFinish(true,xhr.responseText)
	               }
	           }else{
				 	if(xhr.readyState==4){
						option.onFinish(false,"upload failed")
					}     
	           }
	       	   }
			}
			if(option.onProgress){		  
		       xhr.upload.onprogress = function(event){
		           option.onProgress(100 * event.loaded / event.total+"%");
		       }
			}
	      	xhr.send(fd);
			
		  }
	};
	

});