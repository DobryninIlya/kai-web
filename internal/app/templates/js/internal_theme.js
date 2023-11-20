let tg = window.Telegram.WebApp;
if (tg.themeParams != {}) {
    console.log(tg.themeParams.text_color)
    console.log(tg.themeParams)
    document.documentElement.style.setProperty('--background-color', tg.themeParams.bg_color);
    document.documentElement.style.setProperty('--text-color', tg.themeParams.text_color);
    document.documentElement.style.setProperty('--text-hint-color', tg.themeParams.hint_color);
    document.documentElement.style.setProperty('--button-color', tg.themeParams.button_color);
    // document.documentElement.style.setProperty('--button-text-color', tg.themeParams.button_text_colorString);

}