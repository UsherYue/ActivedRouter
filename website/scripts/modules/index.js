//默认的首页
var currentPageID="indexcontent";
//router info初始化的时候会加载
var routerInfo=null ;
var templateEngine =null;
//cpu
var cpuUsedPercent=0;
var cpuFreePercent=0;

//字节到gb
function bytesToGB(bytes){
	var gb=(bytes/1024/1024/1024).toFixed(2).toString();
	return gb ;
}

//加载script模板
var loadScriptTpl=function(scriptid,data){	
	var src=$("#"+scriptid).get(0).src;
	var tpl="";
	$.get(src,function(data){
		tpl=data;
	});
	var render=templateEngine.compile(tpl);
	//返回渲染函数
	var html=render(data);
	return html;
};

//内存初始化
var initMemPie=function(){
	var ctx = $("#memPie").get(0).getContext("2d");
	// For a pie chart
	var usedPercent=routerInfo.VM.used_percent.toFixed(1);
	var freePercent=(100-usedPercent).toFixed(1);
	var data = {
	    labels: [
	        "Mem Used ({0}%)".format(usedPercent),
	        "Mem Free  ({0}%)".format(freePercent),
	
	    ],
	    datasets: [
	        {
	            data: [usedPercent, freePercent],
	            backgroundColor: [
	                "#FF6384",
	                "#36A2EB"
	           
	            ],
	            hoverBackgroundColor: [
	                "#FF6384",
	                "#36A2EB"   
	            ]
	    }]
	};
    //绘制饼图
	var myPieChart = new Chart(ctx,{
	    type: 'pie',
	    data: data,
	    options: null
	});
}

//cpu初始化
var initCpuPie=function(){
	var ctx = $("#cpuPie").get(0).getContext("2d");
	// For a pie chart
	var data = {
	    labels: [
	        "CPU Used ({0}%)".format(cpuUsedPercent),
	        "CPU Free ({0}%)".format(cpuFreePercent),
	
	    ],
	    datasets: [
	        {
	            data: [cpuUsedPercent, cpuFreePercent],
	            backgroundColor: [
	                "#FF6384",
	                "#36A2EB"
	           
	            ],
	            hoverBackgroundColor: [
	                "#FF6384",
	                "#36A2EB"
	          
	            ]
	    }]
	};
    //绘制饼图
	var myPieChart = new Chart(ctx,{
	    type: 'pie',
	    data: data,
	    options: null
	});
	console.log(routerInfo);
}

//cpu初始化
var initDiskPie=function(){
	var ctx = $("#diskPie").get(0).getContext("2d");
	var usedPercent=routerInfo.DISK.used_percent.toFixed(1);
	var freePercent=(100-usedPercent).toFixed(1);
	// For a pie chart
	var data = {
	    labels: [
	        "Disk Used  ({0}%)".format(usedPercent),
	        "Disk Free  ({0}%)".format(freePercent),
	    ],
	    datasets: [
	        {
	            data: [usedPercent,freePercent],
	            backgroundColor: [
	                "#FF6384",
	                "#36A2EB"
	           
	            ],
	            hoverBackgroundColor: [
	                "#FF6384",
	                "#36A2EB"
	          
	            ]
	    }]
	};
    //绘制饼图
	var myPieChart = new Chart(ctx,{
	    type: 'pie',
	    data: data,
	    options: null
	});
}


//init http requ   line chart 
var initHttpLineChart=function(){
	var ctx = $("#httpLineChart").get(0).getContext("2d");
	var data = {
		labels : ["-35min","-30min","-25min","-20min","-15min","-10min","当前状态"],
		datasets : [
			{   
				label:"UIA集群",
	            //是否填充
	            fill: false,
	            // Tension - bezier curve tension of the line. Set to 0 to draw straight lines connecting points
	            // Used to be called "tension" but was renamed for consistency. The old option name continues to work for compatibility.
	            lineTension: 0.1,
	            // 设置贝塞尔曲线下的argb颜色
	            backgroundColor: "rgba(75,192,192,0.4)",
	            //设置曲线颜色
	            borderColor: "rgba(75,192,192,1)",
				//线帽样式
	            borderCapStyle: 'butt',
	            // Array - Length and spacing of dashes. See https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/setLineDash
	            borderDash: [],
	            // Number - Offset for line dashes. See https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/lineDashOffset
	            borderDashOffset: 0.0,
	            // String - line join style. See https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/lineJoin
	            borderJoinStyle: 'miter',
	            // The properties below allow an array to be specified to change the value of the item at the given index
	            // String or Array - Point stroke color
	            pointBorderColor: "rgba(75,192,192,1)",
	            // String or Array - Point fill color
	            pointBackgroundColor: "#fff",
	            // Number or Array - Stroke width of point border
	            pointBorderWidth: 1,
	            // Number or Array - Radius of point when hovered
	            pointHoverRadius: 5,
	            // String or Array - point background color when hovered
	            pointHoverBackgroundColor: "rgba(75,192,192,1)",
	            // String or Array - Point border color when hovered
	            pointHoverBorderColor: "rgba(220,220,220,1)",
	            // 点的圆角 当 hover的时候
	            pointRadius: 1,
	            //数据
	            data: [65, 59, 80, 81, 56, 55, 122],
	            // String - If specified, binds the dataset to a certain y-axis. If not specified, the first y-axis is used. First id is y-axis-0
	            yAxisID: "y-axis-0",
			},
			{   
				label:"学情分析",
				fill: false,
	            backgroundColor: "rgba(255,205,86,0.4)",
	            borderColor: "rgba(255,205,86,1)",
	            pointBorderColor: "rgba(255,205,86,1)",
	            pointBackgroundColor: "#fff",
	            pointBorderWidth: 1,
	            pointHoverRadius: 5,
	            pointHoverBackgroundColor: "rgba(255,205,86,1)",
	            pointHoverBorderColor: "rgba(255,205,86,1)",
	            pointHoverBorderWidth: 2,
	            data: [28, 48, 40, 19, 86, 27, 190]
			},
			{   
				label:"统一资源平台",
				fill: false,
	            backgroundColor: "rgba(175,192,192,0.4)",
	           borderColor: "rgba(99,152,192,1)",
	            pointBorderColor: "rgba(222,66,222,1)",
	            pointBackgroundColor: "#fff",
	            pointBorderWidth: 1,
	            pointHoverRadius: 5,
	            pointHoverBackgroundColor: "rgba(255,205,86,1)",
	            pointHoverBorderColor: "rgba(255,205,86,1)",
	            pointHoverBorderWidth: 2,
	            data:  [25, 47, 36, 19, 45, 66, 121]
			}
		]
	};
	//myLineChart
	var myLineChart = new Chart(ctx,{
	    type: 'line',
	    data: data,
	    options: null
	});
}




//加载content
var loadIndexContent=function(){
	//加载路由服务器信息
	$.get("/routerinfo",function(data){
		routerInfo=data;
		var cpuPercentArr=routerInfo.CPUPERCENTS;
		var usedPercent=0;
		for (var i=0;i<cpuPercentArr.length;i++){
			usedPercent+=cpuPercentArr[i];
		}
		cpuUsedPercent=(usedPercent/cpuPercentArr.length).toFixed(1);
		cpuFreePercent=(100-cpuUsedPercent).toFixed(1)
		routerInfo.cpuFree=cpuFreePercent ;
		routerInfo.cpuUsed=cpuUsedPercent
	});
	//加载模板
    var html=loadScriptTpl("tpl_indexconent",routerInfo);
	//设置html
	$("#center_content").html(html);
	//初始化内存
	initMemPie();
	//初始化cpu
	initCpuPie();
	//初始化磁盘
	initDiskPie();
	//init http
	initHttpLineChart();

};


//加载footer
var loadIndexFooter=function(){
	$('#footer').load('tpl/index_footer.html');  
};


//导出模块
var indexModule=function($,template,Chart,Tools){
	templateEngine=template;
	//初始化格式化函数
	Tools.StringFormatInit();
	//同步ajax
	$.ajaxSetup({  
   		 async : false  
	});  
	//加载首页
	loadIndexContent();
	//加载页脚
	loadIndexFooter();
	//加载活跃主机
	$("#activehost").click(function(){
		//加载模板
    	var html=loadScriptTpl("tpl_indexactive",routerInfo);
		$("#center_content").html(html);
		$('#activehost').parent().siblings().removeClass('active');
		$('#activehost').parent().removeClass('active').addClass('active');
		//$("#clientinfos tbody tr:even").addClass("success");
		$("#clientinfos tbody tr:odd").addClass("warning");
		$("#clientinfos tbody tr").css("height","40px");
		$("#clientinfos tbody tr").hover(function(){
			$(this).css("cursor","pointer");
			console.log("hoverd!");
		});
		$("#clientinfos tbody tr button").click(function(){
			alert("查看服务器详细情况.....")
		});
	});
	//加载index
	$("#indexcontent").click(function(){
		loadIndexContent();
		$('#indexcontent').parent().siblings().removeClass('active');
		$('#indexcontent').parent().removeClass('active').addClass('active');
	});

}

//定义模块
define(['jquery','template','chartjs','tools'],indexModule);