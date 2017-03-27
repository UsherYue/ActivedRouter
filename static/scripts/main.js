(function(){
		require.config({
		//baseUrl:"libs",配置baseurl
		paths:{
			"jquery":"libs/jquery.min",  //jQuery文件。务必在bootstrap.min.js 之前引入 
			"bootstrap":"libs/bootstrap3.3.5.min", //最新的 Bootstrap 核心 JavaScript 文件
			"chartjs":"libs/chart",
			"template":"libs/arttemplate",
			"index":"modules/index",
			"tools":"modules/tools"
		},
		shim : {  
	        bootstrap : {  
	         deps : [ 'jquery' ],  
	         exports : 'bootstrap' 
			},
			chartjs:{
				exports:'Chart'
			}
	    }  
	});
	//引用相关模块 公共模块在此处引入
	requirejs(['bootstrap','index'], function() {  
	
	});  
})();