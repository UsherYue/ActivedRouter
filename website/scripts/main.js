(function(){
		require.config({
		//baseUrl:"libs",配置baseurl
		paths:{
			"jquery":"libs/jquery.min",  //jQuery文件。务必在bootstrap.min.js 之前引入 
			"bootstrap":"libs/bootstrap3.3.5.min", //最新的 Bootstrap 核心 JavaScript 文件
			"template":"libs/arttemplate",
			"index":"modules/index"
		},
		shim : {  
	        bootstrap : {  
	         deps : [ 'jquery' ],  
	         exports : 'bootstrap' 
			},
	    }  
	});
	//引用相关模块
	requirejs(['bootstrap','index','template'], function() {  
	  // alert('loading...');
	});  
})();