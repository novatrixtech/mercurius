function getParameterByName(name, url) {
    if (!url) url = window.location.href;
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
        results = regex.exec(url);
    if (!results) return undefined;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}

$(window).ready(function () {
    var marca = getParameterByName('marca');
    var modelo = getParameterByName('modelo');
    var combustivel = getParameterByName('combustivel');
    var propriedade = getParameterByName('propriedade');
    var ano = getParameterByName('ano');
    var dataInicio = getParameterByName('dataInicio');
    var dataFim = getParameterByName('dataFim');

    if (marca != undefined && marca != '') {
        $('#marca').find('option')
            .filter(function (i, e) {
                return $(e).text() == marca
            }).attr("selected", "selected");
    }

    if (modelo != undefined && modelo != '') {
        $('#modelo').find('option')
            .filter(function (i, e) {
                return $(e).text() == modelo
            }).attr("selected", "selected");
    }

    if (combustivel != undefined && combustivel != '') {
        $('#combustivel').find('option')
            .filter(function (i, e) {
                return $(e).text() == combustivel
            }).attr("selected", "selected");
    }

    if (propriedade != undefined && propriedade != '') {
        $('#propriedade').find('option')
            .filter(function (i, e) {
                return $(e).text() == propriedade
            }).attr("selected", "selected");
    }

    if (ano != undefined && propriedade != '') {
        $('#ano').find('option')
            .filter(function (i, e) {
                return $(e).text() == ano
            }).attr("selected", "selected");
    }

    if (dataInicio != undefined && dataInicio != '') {
        $('#data-inicio').val(dataInicio);
    }

    if (dataFim != undefined && dataFim != '') {
        $('#data-fim').val(dataFim)
    }
});