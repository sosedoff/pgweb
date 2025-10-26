function loadSettings() {
    $.getJSON('settings.json', function(data) {
        settings = data;
        applySettings();
    });
}

function saveSettings() {
    $.ajax({
        url: 'settings.json',
        method: 'PUT',
        contentType: 'application/json',
        data: JSON.stringify(settings),
        success: function() {
            console.log('Settings saved successfully.');
        },
        error: function() {
            console.error('Error saving settings.');
        }
    });
}

function applySettings() {
    document.body.style.backgroundColor = settings.darkTheme ? '#333' : '#fff';
    document.querySelector('.navbar').style.backgroundColor = settings.navbarColor;
    document.querySelector('#logo').src = settings.logo;
}

// $(document).ready(function() {
//     loadSettings();
//     // Additional initialization code...
// });
