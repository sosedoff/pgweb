// This file manages the application settings, including loading and saving settings such as navbar color, logo, and theme preferences.

// const settingsFilePath = path.join(__dirname, '../../settings.json');

let appSettings = {
    navbarColor: '#79589f',
    logo: 'default-logo.png',
    darkTheme: false
};

// function loadSettings() {
//     if (fs.existsSync(settingsFilePath)) {
//         const data = fs.readFileSync(settingsFilePath);
//         appSettings = JSON.parse(data);
//     }
// }
// function loadSettings() {
//     $.getJSON('settings.json', function(data) {
//         appSettings = data;
//         applySettings();
//     });
// }

// function saveSettings() {
//     fs.writeFileSync(settingsFilePath, JSON.stringify(appSettings, null, 2));
// }
function saveSettings() {
    $.ajax({
        url: 'settings.json',
        method: 'PUT',
        contentType: 'application/json',
        data: JSON.stringify(appSettings),
        success: function() {
            console.log('Settings saved successfully.');
        },
        error: function() {
            console.error('Error saving settings.');
        }
    });
}

function applySettings(settings) {
    // document.body.style.backgroundColor = settings.darkTheme ? '#333' : '#fff';
    // document.querySelector('.navbar').style.backgroundColor = appSettings.navbarColor;
    document.documentElement.style.setProperty('--navbar-color', settings.navbarColor);
    console.log(settings)
    $('#sidebar .sidebar-brand > img').attr('src', settings.logo);
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

// function getSettings() {
//     return appSettings;
// }

// Load settings on startup
// loadSettings();

// module.exports = {
//     setNavbarColor,
//     setLogo,
//     toggleTheme,
//     getSettings,
//     loadSettings
// };