let appSettings = {
    navbarColor: '#79589f',
    logo: false,
    darkTheme: false
};

function saveSettings() {
    $.ajax({
        url: 'settings.json',
        method: 'PUT',
        contentType: 'application/json',
        data: JSON.stringify(appSettings),
        success: function () {
            console.log('Settings saved successfully.');
        },
        error: function () {
            console.error('Error saving settings.');
        }
    });
}

function applySettings(settings) {
    document.documentElement.style.setProperty('--navbar-color', settings.navbarColor);
    if (settings.logo) {
        // Display logo wrapper
        $('#sidebar .sidebar-brand').css('display', 'block');
    }
}

function setNavbarColor(color) {
    appSettings.navbarColor = color;
    saveSettings();
}

function setLogo(logo) {
    appSettings.logo = logo;
    saveSettings();
}

function toggleTheme() {
    appSettings.darkTheme = !appSettings.darkTheme;
    saveSettings();
}
