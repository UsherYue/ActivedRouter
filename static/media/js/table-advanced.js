var TableAdvanced = function () {
	
	var _ActiveHostData=null ;
	var loadActiveHost=function(){
		
		
		$.ajax({
			method:"get",
			cache:false ,
			url:"/clientinfos",
			success:function(data){	
				console.log(data);				
			}
		});

		$.get("/clientinfos",function(data){	
			console.log(data);
			_ActiveHostData=data;
	    	    var templateContent=$("#activeListTemplate").html();
			var templateEngine=_.template(templateContent) ;
			var htmlContent=templateEngine({"data":data});
			$("#activeList").html(htmlContent);
			initActiveList();  
		});

	};
		
	//初始化
    var initActiveList = function() {
        function fnFormatDetails ( oTable, nTr )
        {
            var aData = oTable.fnGetData( nTr );
			//获取ip
			var item=_ActiveHostData[aData[2]];
			console.log(item);
            var sOut = '<table>';
            sOut += '<tr><td>服务器IP:</td><td>'+ aData[2]+'</td></tr>';
			sOut += '<tr><td>服务器内核:</td><td>'+ item.Info.HOST.platform+'</td></tr>';
			sOut += '<tr><td>CPU核数:</td><td>'+ item.Info.CpuNums+'</td></tr>';
			var cpus="";
			var allPercent=0;
			for(var i=0;i<item.Info.CPUPERCENTS.length;i++){
				cpus+='<span "color:red;">'+item.Info.CPUPERCENTS[i].toFixed(1).toString()+"</sapn>   ";
				allPercent+=item.Info.CPUPERCENTS[i];
			}
			sOut += '<tr><td>CPU使用百分比/核心:</td><td>'+ cpus+'</td></tr>';
			sOut += '<tr><td>CPU百分比:</td><td>'+ (allPercent/item.Info.CPUPERCENTS.length).toFixed(2).toString()+'</td></tr>';
			sOut += '<tr><td>内核版本:</td><td>'+  item.Info.HOST.platform_version+'</td></tr>';
            sOut += '<tr><td>服务器集群:</td><td>'+item.Info.Cluster+'</td></tr>';
            sOut += '<tr><td>服务器域名:</td><td>'+item.Info.Domain+'</td></tr>';
			sOut += '<tr><td ><strong style="color:red">服务器磁盘信息:</strong></td></tr>';
           	sOut += '<tr><td>文件系统类型:</td><td>'+item.Info.DISK.fstype+'</td></tr>';
			sOut += '<tr><td>服务器存储:</td><td>'+(item.Info.DISK.total/(1024*1024*1024)).toFixed(2).toString()+'GB</td></tr>';
			sOut += '<tr><td>空闲存储:</td><td>'+(item.Info.DISK.free/(1024*1024*1024)).toFixed(2).toString()+'GB</td></tr>';
			sOut += '<tr><td>已经使用存储:</td><td>'+(item.Info.DISK.used/(1024*1024*1024)).toFixed(2).toString()+'GB</td></tr>';
			sOut += '<tr><td>使用百分比:</td><td>'+item.Info.DISK.used_percent.toFixed(2).toString()+'%</td></tr>';
			sOut += '<tr><td ><strong style="color:red">服务器负载状态:</strong></td></tr>';
			sOut += '<tr><td>Last 1 mins:</td><td>'+item.Info.LD.load1.toFixed(2).toString()+'</td></tr>';
           	sOut += '<tr><td>Last 5 mins:</td><td>'+item.Info.LD.load5.toFixed(2).toString()+'</td></tr>';
		    sOut += '<tr><td>Last 15 mins:</td><td>'+item.Info.LD.load15.toFixed(2).toString()+'</td></tr>';
			sOut += '<tr><td ><strong style="color:red">服务器内存:</strong></td></tr>';
			sOut += '<tr><td>服务器内存:</td><td>'+(item.Info.VM.total/1024/1024/1024).toFixed(2).toString()+'GB</td></tr>';
           	sOut += '<tr><td>已使用百分比:</td><td>'+item.Info.VM.used_percent+'%</td></tr>';
			sOut += '<tr><td>已经使用内存:</td><td>'+(item.Info.VM.used/1024/1024/1024).toFixed(2).toString()+'GB</td></tr>';
		    sOut += '<tr><td>空闲内存:</td><td>'+(item.Info.VM.free/1024/1024/1024).toFixed(2).toString()+'GB</td></tr>';
			sOut += '<tr><td>Inactive:</td><td>'+(item.Info.VM.inactive/1024/1024/1024).toFixed(2).toString()+'GB</td></tr>';
			sOut += '<tr><td ><strong style="color:red">服务器TCP连接网络:</strong></td></tr>';
			sOut += '<tr><td>所有连接数:</td><td>'+item.Info.NC.ALLCOUNT+'</td></tr>';
			sOut += '<tr><td>CLOSED_WAIT:</td><td>'+item.Info.NC.CLOSED_WAIT+'</td></tr>';
			sOut += '<tr><td>ESTABLISH:</td><td>'+item.Info.NC.ESTABLISH+'</td></tr>';
			sOut += '<tr><td>LISTEN:</td><td>'+item.Info.NC.LISTEN+'</td></tr>';
			sOut += '<tr><td>TIME_WAIT:</td><td>'+item.Info.NC.TIME_WAIT+'</td></tr>';
			sOut += '<tr><td>SYN_SENT:</td><td>'+item.Info.NC.SYN_SENT+'</td></tr>';
			sOut += '<tr><td>SYN_RECV:</td><td>'+item.Info.NC.SYN_RECV+'</td></tr>';
			sOut += '<tr><td>FIN_WAIT_1:</td><td>'+item.Info.NC.FIN_WAIT_1+'</td></tr>';
			sOut += '<tr><td>FIN_WAIT_2:</td><td>'+item.Info.NC.FIN_WAIT_2+'</td></tr>';
			sOut += '</table>';     
            return sOut;
        }
	

        var nCloneTh = document.createElement( 'th' );
        var nCloneTd = document.createElement( 'td' );
        nCloneTd.innerHTML = '<span class="row-details row-details-close"></span>';         
        $('#sample_1 thead tr').each( function () {
            this.insertBefore( nCloneTh, this.childNodes[0] );
        } );
         
        $('#sample_1 tbody tr').each( function () {
            this.insertBefore(  nCloneTd.cloneNode( true ), this.childNodes[0] );
        } );
        var oTable = $('#sample_1').dataTable( {
            "aoColumnDefs": [
                {"bSortable": false, "aTargets": [ 0 ] }
            ],
            "aaSorting": [[1, 'asc']],
             "aLengthMenu": [
                [5, 15, 20, -1],
                [5, 15, 20, "All"] // change per page values here
            ],
            // set the initial value
            "iDisplayLength": 10,
        });

        jQuery('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap small"); // modify table search input
        jQuery('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
        jQuery('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown
        $('#sample_1').on('click', ' tbody td .row-details', function () {
            var nTr = $(this).parents('tr')[0];
            if ( oTable.fnIsOpen(nTr) )
            {
                $(this).addClass("row-details-close").removeClass("row-details-open");
                oTable.fnClose( nTr );
            }
            else
            {
                $(this).addClass("row-details-open").removeClass("row-details-close");
                oTable.fnOpen( nTr, fnFormatDetails(oTable, nTr), 'details' );
            }
        });
    }
    //导出
    return {
        init: function () {
            if (!jQuery().dataTable) {
                return;
            }
			loadActiveHost();       
        },
		reinit:function(){
			loadActiveHost();
		}
    };

}();