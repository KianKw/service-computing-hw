$(document).ready(function() {
    $.ajax({
        url: "/api/test"
    }).then(function(data) {
       $('.name').append(data.name);
       $('.date').append(data.date);
    });
});
