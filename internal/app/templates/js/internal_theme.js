let tg = window.Telegram.WebApp;
const root = document.documentElement;
const colorScheme = getComputedStyle(root).getPropertyValue('--tg-color-scheme');
if (colorScheme != 'light' && colorScheme != '') {
    console.log(tg.themeParams.text_color)
    console.log(tg.themeParams)
    root.style.setProperty('--background-color', tg.themeParams.bg_color);
    root.style.setProperty('--text-color', tg.themeParams.text_color);
    root.style.setProperty('--text-hint-color', tg.themeParams.hint_color);
    root.style.setProperty('--text-color-depp-gray', tg.themeParams.hint_color);
    root.style.setProperty('--button-color', tg.themeParams.button_color);
    root.style.setProperty('--hashtag-color', tg.themeParams.hint_color);
    root.style.setProperty('--secondary-background-color', tg.themeParams.secondary_bg_color);
    // document.documentElement.style.setProperty('--button-text-color', tg.themeParams.button_text_colorString);
    const arrowElements = document.querySelectorAll('.arrow');
    arrowElements.forEach((arrow) => {
        arrow.setAttribute('fill', tg.themeParams.text_color);
    });
    tg.expand()
}