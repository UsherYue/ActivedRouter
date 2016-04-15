#   ActiveRouter简介    　
当前版本: Beta 0.1	by usher.yue	

Email:   usher.yue@gmail.com		
			
<br>
##  <b>简单介绍</b>
	1、ActiveRouter分为Sever 模式、Client模式以Reserve Proxy 模式 。通过ActiveRouter
	可以实现内网环境下的api负载均衡、服务器监控等功能。
        
	2、提供服务器运行状态监控功能,可时刻返回服务器状态,对服务器运行状态了如指掌。
        
	3、后期提供第三方SDK php、golang、java、nodejs、ruby等等,实现http请求的真机负载。
        
	4、Server & Client 模式组合下,可以用作服务器监控使用,并且可以设置监控报警,Server端可以选择性开启监控Dashboard,
	用于直观的监控。
        
	5、Proxy模式下可以提供简单的反向代理服务。 
        
	6、=_=!!!  所有以上功能都在完善中～～～～
##  快速入门		
###  编译安装
在系统下运行如下脚本,会在当前目录

	linux64下    boot_linux64.sh
	darwin64下   boot_darwin64.sh
	linux64下    boot_linux64.sh
               
##  工作模式 
### 1、反向代理模式 Reserve Proxy ,类似nginx的反向代理功能
运行命令： ActivedRouter --runmode=proxy		

配置文件：
        
 		{
			"proxy_addr":"127.0.0.1:80",
			"proxy_method":"random",
			"reserve_proxy":[
				{
					"domain":"api.golang.com",
					"clients":[
					    {
					 	  "host":"uia.com",
						   "port":"80"	
		    	   		 },
						{
					 	  "host":"uiatest.com",
					 	  "port":"80"	
		    	    	}
					]
				},
				{
					"domain":"api1.golang.com",
					"clients":[
					    {
						   "host":"goapi.com",
						   "port":"80"	
		      		  }
					]
				}
			]
		}
### 2、server和client模式可以配合完全服务器监控,以及路由转发
        ActivedRouter --runmode=server  
        
        ActivedRouter --runmode=client  
### 3、钩子脚本 
	{   
			"script":[
		{
				"host":"127.0.0.1",
				"hookscript":[
					{
						"eventtarget":"disk",
						"attr":"used",
						"gt":"15",
						"callback":"ls"
					},
					{
						"eventtarget":"load",
						"attr":"load",
						"gt":"1.5",
						"callback":"ls"
					},
					{
						"eventtarget":"mem",
						"attr":"used",
						"gt":"75",
						"callback":"netstat -ant"
					}
				]
			}
		]
	}
##  <b>运行模式</b>
        服务器模式和客户端模式
<table >
   <thead>
     <tr>
        <th>运行模式</th>
        <th>介绍</th>
     </tr>
   </thead>
   <tbody>
    <tr>
      <td>
         ServerMode
      </td>
      <td>
            在服务器模式下监听客户端的状态 <br/>
            1、 第三方应用通过sdk提供的功能进行自动路由到合适的api服务器上,此处属于代理转发功能。<br/>
            2、 第三方应用通过sdk提供的功能获取到一个空闲合适的服务器域名或ip。<br/>
            3、 监听客户端模式下的服务器的服务状态
      </td>
    </tr>
    <tr>
      <td>
         ClientMode
      </td>
      <td>
            在客服端模式下通告服务器当前服务状态 <br/>
            1、 启动挂在到路由服务器 可以挂在到一个至个服务器上。<br/>
            2、 实时通告路由服务器当前服务器信息,用作路由分析。
      </td>
    </tr>
	<tr>
      <td>
         Reserve Proxy
      </td>
      <td>
         内网集群配置反向代理功能
      </td>
    </tr>
    <tr>
      <td>
         ThirdPartSDK
      </td>
      <td>
        针对第三方提供php golang 等sdk,提供基于路由负载、以及反向代理负载的http请求机制。
      </td>
    </tr>
   </tbody>
</table>    
##  提供api服务器监控功能可以实时返回各服务器状态  
<table >
   <thead>
     <tr>
        <th>监控功能</th>
        <th>介绍</th>
     </tr>
   </thead>
   <tbody>
    <tr>
      <td>
         虚拟内存
      </td>
      <td>
             时刻监控服务器的虚拟内存
      </td>
    </tr>
     <tr>
      <td>
         load average
      </td>
      <td>
             时刻监控服务器的负载状态
      </td>
    </tr>
    <tr>
      <td>
         网络连接
      </td>
      <td>
        时刻监控服务器的网络连接状态
      </td>
    </tr>
    <tr>
      <td>
         磁盘状态
      </td>
      <td>
        时刻监控服务器的磁盘存储容量
      </td>
    </tr>
    <tr>
      <td>
         ThirdPartSDK
      </td>
      <td>
         为应用提供监控接口,可直接展示监控内容
      </td>
    </tr>
   </tbody>
</table>