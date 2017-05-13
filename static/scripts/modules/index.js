//默认的首页
var currentPageID="indexcontent";
//router info初始化的时候会加载
var routerInfo=null ;
//cpu
var cpuUsedPercent=0;
var cpuFreePercent=0;
var hostClientsInfo=null;

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

//init http request   line chart 
var initHttpLineChart=function(){
	var colorTable=[
		["rgba(75,192,192,0.4)","rgba(75,192,192,1)"],
		["rgba(225,192,92,0.4)","rgba(225,192,92,1)"],
		["rgba(175,92,192,0.4)","rgba(175,92,192,1)"],	
		["rgba(65,192,92,0.4)","rgba(65,192,92,1)"],
		["rgba(115,32,92,0.4)","rgba(115,32,92,1)"],
		["rgba(65,192,92,0.4)","rgba(65,192,92,1)"],
		["rgba(115,32,92,0.4)","rgba(115,32,92,1)"]
	];
	//加载路由服务器信息
	$.get("/statistics",function(data){
		// alert(JSON.stringify(data));
		var colorIndex=0;
		  //lables
		  var labels=['-50','-45min','-40min','-35min','-30min','-25min','-20min','-15min','-10min','-5min','当前'];
		 //每条曲线的数据集合
		  var datasets=[];
			for(var key in data){
				//每一个集群
				var tmpData=[];
				var tmpLen=data[key].length;
				var endIndex=tmpLen-1-10;
				var timeSep=0;
				var tmpReqCount=0;
				for(var i=endIndex;i<=tmpLen-1;i++ ){
					if(i<0){
						tmpReqCount=0;
					}else{
						tmpReqCount=data[key][i].RequestCount;
					}
					timeSep=tmpLen-1-i;
					tmpData.push(tmpReqCount)
				}
				datasets.push({   
					label:key,
		            //是否填充
		            fill: false,
		            // 设置贝塞尔曲线下的argb颜色
		            backgroundColor: colorTable[colorIndex%colorTable.length][0],
		            //设置曲线颜色
		            borderColor: colorTable[colorIndex%colorTable.length][1],
					 pointBackgroundColor: "#fff",
		            pointBorderWidth: 1,
		            pointHoverRadius: 5,
		            pointHoverBackgroundColor: "rgba(255,205,86,1)",
		            pointHoverBorderColor: "rgba(255,205,86,1)",
		            pointHoverBorderWidth: 2,
		            //数据
		            data: tmpData,
				});
				colorIndex++;
			}
			//渲染chart
			var ctx = $("#httpLineChart").get(0).getContext("2d");
			var data = {
				labels : labels,
				datasets :datasets 
			};
			//myLineChart
			var myLineChart = new Chart(ctx,{
			    type: 'line',
			    data: data,
			    options: null
			});
	});
}

//op=0 删除  op=1编辑修改
function InitProxyInsertRowEvent(op,data){
	    var hostInput=$("#proxy_setting_dlg table tr.insert-row input:eq(0)");
		var portInput=$("#proxy_setting_dlg table tr.insert-row input:eq(1)");
		//确认添加
		if(!op){
			$("#proxy_setting_dlg table tr.insert-row button:eq(0)").click(function(){
			   var host=hostInput.val();
			   var port=portInput.val();
               if(!/^[\w\.]{1,255}$/i .test(host)||!/^[0-9]{1,5}$/i .test(port)){
				  $("#proxy_setting_dlg table tr.insert-row td:eq(3)").html('<div class="bg-danger" >参数不合法</div>');
			   }else{
					var domain=$("#default-domain").text();
					 $.get('/addproxyclient',{host:host,port:port,domain:domain},function(data){
						  $.get("/proxyinfos/{0}".format(domain),function(data){
								var res=templateEngine("Proxy_Client_List",{defaultHosts:data});
								$("#proxy_setting_dlg div.proxy-domain-client").html(res);
								$("#default-domain").text(domain);
								InitProxyRowEvent();
							});
					});
				}	 
			});
		}else{
			//确认修改
			$("#proxy_setting_dlg table tr.insert-row button:eq(0)").click(function(){
			   var host=hostInput.val();
			   var port=portInput.val();
               if(!/^[\w\.]{1,255}$/i .test(host)||!/^[0-9]{1,5}$/i .test(port)){
				  $("#proxy_setting_dlg table tr.insert-row td:eq(3)").html('<div class="bg-danger" >参数不合法</div>');
			   }else{
					var domain=$("#default-domain").text();
					$.get("/updateproxyclient",{domain:domain,prehost:data.host,preport:data.port,updatehost:host,updateport:port},function(data){
						 if(data.status==1){
							selectRow=$("#proxy_setting_dlg table tr.insert-row").prev();
							selectRow.show();
							$("#proxy_setting_dlg table tr.insert-row").remove();
							//更新为修改后的结果
							selectRow.find("td:eq(0)").text(host);
							selectRow.find("td:eq(1)").text(port);
						}else{
							$("#proxy_setting_dlg table tr.insert-row td:eq(3)").html('<div class="bg-danger" >修改失败!</div>');
						}
					});
				}	 
			});
		}
		//取消可输入框	
		$("#proxy_setting_dlg table tr.insert-row button:eq(1)").click(function(){
			//显示编辑隐藏的行
			 if(op==1)$(this).parent().parent().prev().show();
			 hostInput.val('');
			 portInput.val('');
			 hostInput.parent().hasClass('has-success')?hostInput.parent().removeClass('has-success'):false
			 hostInput.parent().hasClass('has-error')?hostInput.parent().removeClass('has-error'):false
			 hostInput.next().hasClass('glyphicon-ok')?hostInput.next().removeClass('glyphicon-ok'):false
		 	 hostInput.next().hasClass('glyphicon-remove')?hostInput.next().removeClass('glyphicon-remove'):false
			 portInput.parent().hasClass('has-success')?portInput.parent().removeClass('has-success'):false
			 portInput.parent().hasClass('has-error')?portInput.parent().removeClass('has-error'):false
			 portInput.next().hasClass('glyphicon-ok')?portInput.next().removeClass('glyphicon-ok'):false
		 	 portInput.next().hasClass('glyphicon-remove')?portInput.next().removeClass('glyphicon-remove'):false
			$("#proxy_setting_dlg table tr.insert-row td:eq(3)").html('');
			$("#proxy_setting_dlg table tr.insert-row").remove();
			
		});
		//输入框控制		
		hostInput.keyup(function(data){
			 var host=$(this).val()
			 if(!/^[\w\.]{1,255}$/i .test(host)){
					$("#proxy_setting_dlg table tr.insert-row td:eq(3)").html('<div class="bg-danger" >ip或者域名不合法</div>');
					hostInput.parent().removeClass('has-success');
					hostInput.parent().addClass('has-error');
					hostInput.next().removeClass('glyphicon-ok');
					hostInput.next().addClass('glyphicon-remove');
			}else{
					hostInput.parent().removeClass('has-error');
					hostInput.next().removeClass('glyphicon-remove');
					hostInput.parent().addClass('has-success');
					hostInput.next().addClass('glyphicon-ok');
					$("#proxy_setting_dlg table tr.insert-row td:eq(3)").html('');
			}	    
		});
		portInput.keyup(function(data){
			 var port=$(this).val()
			 if(!/^[0-9]{1,5}$/i .test(port)){
					$("#proxy_setting_dlg table tr.insert-row td:eq(3)").html('<div class="bg-danger" >端口不合法</div>');
					portInput.parent().removeClass('has-success');
					portInput.parent().addClass('has-error');
					portInput.next().removeClass('glyphicon-ok');
					portInput.next().addClass('glyphicon-remove');
			}else{
					portInput.parent().removeClass('has-error');
					portInput.next().removeClass('glyphicon-remove');
					portInput.parent().addClass('has-success');
					portInput.next().addClass('glyphicon-ok');
					$("#proxy_setting_dlg table tr.insert-row td:eq(3)").html('');
			}	    
		})
}

//初始化反向代理列表普通行的事件
function InitReservProxyClientEvent(){
	//删除域名
	$("#proxy_setting_dlg table tr td a.remove-client").click(function(){
		var domain=$("#default-domain").text();
		var host=$(this).parent().parent().children("td").eq(0).text();
		var port=$(this).parent().parent().children("td").eq(1).text();
		 $.get('/delproxyclient',{host:host,port:port,domain:domain},function(data){
			  $.get("/proxyinfos/{0}".format(domain),function(data){
					var res=templateEngine("Proxy_Client_List",{defaultHosts:data});
					$("#proxy_setting_dlg div.proxy-domain-client").html(res);
					$("#default-domain").text(domain);
					//重新初始化ProxyRowEvent
					InitProxyRowEvent();
				});
		});  
	});
	//编辑域名
	$("#proxy_setting_dlg table tr td a.edit-client").click(function(){
		//删除前一个编辑框
		var insertObj=$("#proxy_setting_dlg table tr.insert-row");
		if(insertObj.length>0){
			insertObj.prev().show();
			insertObj.remove();
		}
		var html=templateEngine("Proxy_Client_EditRow",{});
		//将edit 插入到选定行下一行,并且隐藏本行
		var selectRow=$(this).parent().parent();
		selectRow.after(html);
		//隐藏选择航
		selectRow.hide();
		//插入edit row
		insertObj=selectRow.next();
		var hostInput=insertObj.find("input:eq(0)");
		var portInput=insertObj.find("input:eq(1)");
		var defaultDomain=$("#default-domain").text();
		var preHost=selectRow.find("td:eq(0)").text();
		var prePort=selectRow.find("td:eq(1)").text();
		hostInput.val(preHost);
		portInput.val(prePort);
		//初始化编辑框事件
		InitProxyInsertRowEvent(1,{host:preHost,port:prePort});
	});
}

function ChangeInputStatus(jqObj,bSuccess){
	if(!bSuccess){
		jqObj.parent().removeClass('has-success');
		jqObj.parent().addClass('has-error');
		jqObj.next().removeClass('glyphicon-ok');
		jqObj.next().addClass('glyphicon-remove');
	}else{
		jqObj.parent().removeClass('has-error');
		jqObj.next().removeClass('glyphicon-remove');
		jqObj.parent().addClass('has-success');
		jqObj.next().addClass('glyphicon-ok');
	}
}

function InitDomainEditEvent(op){
	//增加新域名如果域名存在则增加失败
	$("#domain_setting_dlg div.btn-add-domain").click(function(){
		var insertRow=$("#domain_setting_dlg table tr[flag='domain-edit']");
		//保证只出现一个插入框
		if(insertRow.length==0){
		   var insertHtml=templateEngine("Proxy_Domain_EditRow",{});
		  //行尾插入
		   $("#domain_setting_dlg table ").append(insertHtml);
		   insertRow=$("#domain_setting_dlg table tr[flag='domain-edit']");
		   insertRow.find("button:eq(0)").click(function(){
				//添加域名
				var addDomain=insertRow.find("input").val();
				if(!/^[\w\.]{1,255}$/i .test(addDomain)){
					alert('域名不合法!...');				
				}else{
					$.get("/adddomain/{0}".format(addDomain),function(data){
						//添加完成 
						$("#domain_setting_dlg table tr:last").after("<tr><td>{0}</td><td><a href=\"#\" class=\"remove-client\">remove</a>&nbsp;/&nbsp;<a href=\"#\" class=\"edit-client\">Edit</a></td></tr>".format(addDomain));
						$("#domain_setting_dlg table tr.insert-row").remove();
					});	
				}
		   });
			$("#inputDomain").keyup(function(data){
			    var  domainVal=$(this).val();
				ChangeInputStatus($(this),/^[\w\.]{1,255}$/i .test(domainVal));
			});

		 $("#domain_setting_dlg table tr.insert-row button:eq(1)").click(function(){
			   $("#domain_setting_dlg table tr.insert-row").remove();
		  });
		}  
	});
	//删除
	$("#domain_setting_dlg table tr td a.remove-domain").click(function(){
			//单击确定删除
		    var bToDo = window.confirm("单击\"确定\"删除。单击\"取消\"停止。");
			if(bToDo){
				var selectRow=$(this).parent().parent();
				var delDomain=$(this).parent().prev().text();
				//删除域名
				$.get("/deldomain/{0}".format(delDomain),function(){
					selectRow.remove();
				});						
			}
	});
	//编辑
	$("#domain_setting_dlg table tr td a.edit-domain").click(function(){
		var insertRow=$("#domain_setting_dlg table tr.insert-row");
		var selectRow=$(this).parent().parent();
		var delDomain=$(this).parent().prev().text();
		//删除上一个编辑框
		if(insertRow.length>0){
			insertRow.prev().show();
			insertRow.remove();	
			selectRow.after(insertRow);
		}else{
			var insertHtml=templateEngine("Proxy_Domain_EditRow",{});
			selectRow.after(insertHtml);
		}
		selectRow.hide();
		$("#inputDomain").val(delDomain);
		//update
		$("#domain_setting_dlg table tr.insert-row button:eq(0)").click(function(){
			var updateDomain=$("#inputDomain").val();
			if(!/^[\w\.]{1,255}$/i .test(updateDomain)){
					alert('域名不合法!...');				
				}else{
					//更新域名
					$.get("/updatedomain",{predomain:delDomain,updatedomain:updateDomain},function(data){
					selectRow.show();
					$("#domain_setting_dlg table tr.insert-row").remove();
					selectRow.find("td:eq(0)").text(updateDomain);			
				});
			}
		});
		$("#inputDomain").keyup(function(data){
		    var  domainVal=$(this).val();
			ChangeInputStatus($(this),/^[\w\.]{1,255}$/i .test(domainVal));
		});
		//remove
		$("#domain_setting_dlg table tr.insert-row button:eq(1)").click(function(){
			selectRow.show();
			$("#domain_setting_dlg table tr.insert-row").remove();
		});			
	});

	
}

//初始化菜单事件
var initMenuEvent=function(){
	//点击加载活跃主机
	$("#activehost").click(function(){
	    loadActiveContent();
		$('#activehost').parent().siblings().removeClass('active');
		$('#activehost').parent().removeClass('active').addClass('active');
	});
	//点击加载index
	$("#indexcontent").click(function(){
		loadIndexContent();
		$('#indexcontent').parent().siblings().removeClass('active');
		$('#indexcontent').parent().removeClass('active').addClass('active');
	});
	//反向代理配置
	$("#http_reserveproxy_setting").click(function(){
		var clientList=[];
		var domainList=[];
		var defaultDomainSelect="";
		//域名信息
		$.get("/domaininfos",function(data,status){
			domainList=data;
			defaultDomainSelect=data.length?data[0]:"No Domain"
			if(data.length>0){			
				$.get("/proxyinfos/{0}".format(defaultDomainSelect),function(data){
					clientList=data;
				});
			}
  		});
		var html=templateEngine("tpl_proxysetting",{defaultDomain:defaultDomainSelect,domainlist:domainList,defaultHosts:clientList});
		$(html).appendTo($("body"));
		//域名选择切换事件
		$("#proxy_setting_dlg ul.select-domain-menu li a").click(function(){
			defaultDomainSelect=$(this).text();
			$.get("/proxyinfos/{0}".format(defaultDomainSelect),function(domainClients){
				var res=templateEngine("Proxy_Client_List",{defaultHosts:domainClients});
				$("#proxy_setting_dlg div.proxy-domain-client").html(res);
				$("#default-domain").text(defaultDomainSelect);
					InitReservProxyClientEvent();
			});
		});
		//add proxy
		$("#btnAddProxy").click(function(){
				var html=templateEngine("Proxy_Client_EditRow",{});
				if($("#proxy_setting_dlg table tr:eq(1)").attr('flag')!="proxy-edit"){
					$("#proxy_setting_dlg table tr:eq(0)").after(html);
					//初始化可编辑表格事件
					InitProxyInsertRowEvent(0);
				}
		});
	    InitReservProxyClientEvent();
		//Open Diload
		$("#proxy_setting_dlg").modal("show");
		$("#proxy_setting_dlg").on('hide.bs.modal',function(){
			$(this).remove();
		});
	});
	
	//域名列表
	$("#domain_setting").click(function(){
		$.get("/domaininfos",function(data,status){
			var html=templateEngine("tpl_domainsetting",{domainInfos:data});
			$(html).appendTo($("body"));
			InitDomainEditEvent();
			$("#domain_setting_dlg").modal("show");
			$("#domain_setting_dlg").on('hide.bs.modal',function(){
				$(this).remove();
			});
  		});	
	});
}

//加载活跃主机列表
var loadActiveContent=function(){
	$.get('/clientinfos',function(data){
		//计算cpu
		for(var key in data){
			var cpuPercent=0.0;
			var cpuLen=data[key].Info.CPUPERCENTS.length;
			for(var i =0;i<data[key].Info.CPUPERCENTS.length;i++){
				cpuPercent+=data[key].Info.CPUPERCENTS[i];
			}
			data[key].Info.CPUPERCENT=(cpuPercent/cpuLen).toFixed(2);
		}
		//保存数据
		hostClientsInfo=data;
		var html=loadScriptTpl("tpl_indexactive",{hostlist:data});
		$("#center_content").html(html);
		$("#clientinfos tbody tr:even").addClass("warning");
		//查看服务器详细情况
		$("#clientinfos tbody td").dblclick(function(){
			var ip=$(this).siblings().eq(1).html();
			var hostInfo=hostClientsInfo[ip];
			var html=loadScriptTpl('tpl_activehost',hostInfo.Info)
			$("#dlg_hostinfo").html(html);
			//初始化内存
			initMemPie();
			//初始化cpu
			initCpuPie();
			//初始化磁盘
			initDiskPie();
			$('#dlg').modal('show')	
		});
		//查看服务器详细情况
		$("#clientinfos tbody tr button").click(function(){
			var ip=$(this).parent().siblings().eq(1).html();
			var hostInfo=hostClientsInfo[ip];
			var html=loadScriptTpl('tpl_activehost',hostInfo.Info)
			$("#dlg_hostinfo").html(html);
			//初始化内存
			initMemPie();
			//初始化cpu
			initCpuPie();
			//初始化磁盘
			initDiskPie();
			$('#dlg').modal('show')	
		});
	});
};

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
	if(routerInfo.RunMode=="reserveproxy"){
			initHttpLineChart();
	}else{
		//非反向代理模式下隐藏折线图
		$("#http-proxy-statistics").css("display","none");
	}
};

//加载footer
var loadIndexFooter=function(){
	$('#footer').load('tpl/index_footer.html');  
};

//加载menu
var loadIndexMenu=function(){
	$('#menu_content').load('tpl/index_menu.html');
}

//导出模块
var indexModule=function($,template,Chart,Tools){
	//关闭模板缓存
	template.config.cache=false;
	window.templateEngine=template;
	//初始化jquery
	window.$=$;
	//同步ajax
	window.$.ajaxSetup({  
   		 async : false ,
		 cache:false
	});  
	//加载menu
	loadIndexMenu();
	//加载首页
	loadIndexContent();
	//加载页脚
	loadIndexFooter();
	//初始化菜单事件
	initMenuEvent();
}

//定义模块
define(['jquery','template','chartjs','tools'],indexModule);