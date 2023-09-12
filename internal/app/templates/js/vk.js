// Для секции VK Storage:
// "join_group" - предлагалось ли подписаться на группу
// "add_to_hs" - предлагалось ли добавить на рабочий экран
//

function initVK() {
    console.log(window.vkBridge)
    window.vkBridge.send('VKWebAppInit').then((data) => {
        console.log("init VKWebAppInit")
    })
        .catch((error) => {
            // Обработка события в случае ошибки
            console.log(error);
        });

    var add_to_hs, join_group;
    window.vkBridge.send("VKWebAppStorageGet", {"keys": ["join_group", "add_to_hs"]})
        .then((data) => {
            const keysValues = data.keys.map(item => ({[item.key]: item.value}));
            join_group = keysValues[0]['join_group'];
            add_to_hs = keysValues[1]['add_to_hs'];
        })
        .catch((error) => {
            console.log(error);
        });

    addToHSInvite(add_to_hs);
    joinGroupInvite(join_group);
}

function addToHSInvite(add_to_hs) {
    let is_feature_supported;
    if (add_to_hs == "" || true) {
        window.vkBridge.send("VKWebAppAddToHomeScreenInfo")
            .then((hs_info) => {
                console.log(hs_info["is_feature_supported"])
                is_feature_supported = hs_info["is_feature_supported"]
                console.log(hs_info["is_added_to_home_screen"])

            })
            .catch((error) => {
                // Обработка события в случае ошибки
                console.log(error);
            });
        if (is_feature_supported) {
            window.vkBridge.send("VKWebAppAddToHomeScreen");
            window.vkBridge.send("VKWebAppStorageSet", {"key": "add_to_hs", "value": "offered"})
                .then((data) => {
                    console.log(data)
                })
                .catch((error) => {
                    // Обработка события в случае ошибки
                    console.log(error);
                });
        } else {
            var element = document.getElementById("add_to_homescreen");
            element.parentNode.removeChild(element);
        }
    }
}

function joinGroupInvite(join_group) {
    if (join_group == "") {
        window.vkBridge.send("VKWebAppJoinGroup", {"group_id": 182372147})
            .then((data) => {
            console.log(data)
        })
            .catch((error) => {
                // Обработка события в случае ошибки
                console.log(error);
            });
        window.vkBridge.send("VKWebAppStorageSet", {"key": "join_group", "value": "offered"})
            .then((data) => {
                console.log(data)
            })
            .catch((error) => {
                // Обработка события в случае ошибки
                console.log(error);
            });
    }
}

