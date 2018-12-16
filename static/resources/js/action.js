
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

function notify(target_cgi)
{
    return function() {
        $.get(
            target_cgi,
            function(data , status){
                alert(data.msg)
            }
        )
    }
}


$("#id_submit").click(saveConfigData);