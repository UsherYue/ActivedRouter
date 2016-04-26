//内存初始化
var initMemPie=function(){
	var ctx = $("#memPie").get(0).getContext("2d");
	// For a pie chart
	var data = {
	    labels: [
	        "Mem Used (35%)",
	        "Mem Free  (65%)",
	
	    ],
	    datasets: [
	        {
	            data: [35, 65],
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
	        "CPU Used (25%)",
	        "CPU Free (75%)",
	
	    ],
	    datasets: [
	        {
	            data: [25, 75],
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
var initDiskPie=function(){
	var ctx = $("#diskPie").get(0).getContext("2d");
	// For a pie chart
	var data = {
	    labels: [
	        "Disk Used  (33.3%)",
	        "Disk Free  (66.7%)",
	    ],
	    datasets: [
	        {
	            data: [33.3, 66.7],
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



var indexModule=function($,template,Chart,Tools){
	//初始化格式化函数
	Tools.StringFormatInit();
	//初始化内存
	initMemPie();
	//初始化cpu
	initCpuPie();
	//初始化磁盘
	initDiskPie();
	//init http
	initHttpLineChart();

}

//定义模块
define(['jquery','template','chartjs','tools'],indexModule);