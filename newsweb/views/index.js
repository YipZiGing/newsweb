window.onload=function (ev) {
    $(".dels").click(function () {
        if (!confirm("确认删除吗？")) {
            return false
        }
    })
    $("#select").change(function () {
        $("#form").submit()
    })
}
