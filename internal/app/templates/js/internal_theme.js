let tg = window.Telegram.WebApp;
if (document.documentElement.style.getProperty('--tg-color-scheme' != 'light')) {
    console.log(tg.themeParams.text_color)
    console.log(tg.themeParams)
    document.documentElement.style.setProperty('--background-color', tg.themeParams.bg_color);
    document.documentElement.style.setProperty('--text-color', tg.themeParams.text_color);
    document.documentElement.style.setProperty('--text-hint-color', tg.themeParams.hint_color);
    document.documentElement.style.setProperty('--button-color', tg.themeParams.button_color);
    document.documentElement.style.setProperty('--hashtag-color', tg.themeParams.hint_color);
    document.documentElement.style.setProperty('--secondary-background-color', tg.themeParams.secondary_bg_color);
    // document.documentElement.style.setProperty('--button-text-color', tg.themeParams.button_text_colorString);
    const arrowElements = document.querySelectorAll('.arrow');
    arrowElements.forEach((arrow) => {
        arrow.setAttribute('fill', tg.themeParams.text_color);
    });
    tg.expand()
}