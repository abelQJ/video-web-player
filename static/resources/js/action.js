
function saveConfigData()
{
    $.post("/cgi/setconfig",
       $("#id_config").val(),
       function(data,status){
        alert(data.msg)
       });
}

function getConfigData()
{
    $.get("/cgi/getconfig",
       function(data,status){
        $("#id_config").val(JSON.stringify(data,null,4));
       });
}

$("#id_submit").click(saveConfigData);