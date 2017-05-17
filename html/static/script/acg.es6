// acg
// @author: zhengxiaoyao0716
; (function () {
    "use strict";
    const acg = {
        "background": {
            "doms": document.querySelectorAll(".-acg- .background, .-acg-.background"),
            "images": [{
                url: "./static/image/acg_bgs/49647343_p0.png",
                info: "PixivID: 49647343 | 画师：コ゛りぼて",
            }],  // 图片缓存队列，可空数组，获取到的图片会扩展进入数组。
            // "images": [],
            // "interval": 30,  // 自动切换时间间隔，单位s，不大于0时（0、-1、undefined）不自动切换
            // "provider": "http://pic.api.freejishu.com/v2/?tag=acg",  // 图片供应api，空值表示不获取新图片
        },
        "colors": ["black", "blue", "green", "yellow", "pink", "indigo", "red", "grey",],
    };

    // 随机背景
    if (acg.background.doms) {
        acg.background.images = acg.background.images || [];
        let index = 0;
        // 拉取下一张图片
        const fetchImage = acg.background.provider ? (fn) => {
            base.fetch(acg.background.provider).then((data) => {
                if (fn) {
                    acg.background.images.push({
                        url: data['url'],
                        info: '' || data['url'],
                    });
                    fn(acg.background.images[index]);
                } else {
                    base.preload.image(data['url'], (src) => {
                        acg.background.images.push({
                            url: src,
                            info: '' | src,
                        });
                    });
                }
            });
        } : (fn) => {
            index = 0;
            fn && fn(acg.background.images[index]);
        };
        // 切换背景到下一张
        function nextImage() {
            new Promise((resolve, reject) => {
                function onImageGet(image) {
                    ++index >= acg.background.images.length && fetchImage();
                    resolve(image);
                }
                const image = acg.background.images[index];
                if (!image) {
                    fetchImage(image => onImageGet(image));
                } else {
                    onImageGet(image);
                }
            }).then(image => image && acg.background.doms.forEach(div => {
                div.style.backgroundImage = 'url(' + image.url + ')';
                div.setAttribute("data-acg-bg-info", image.info);
            }));
        };
        nextImage();
        // 定时刷新
        acg.background.interval > 0 && setInterval(nextImage, acg.background.interval * 1000);
    }

    //数据管理器
    const dataManager = (() => {
        // 解析日期，"2017-04-03" => ["Apr 09, 2017", "green"]
        function parseDate(dateStr) {
            const [year, month, date] = dateStr.split('-');
            const d = new Date(year, month - 1, date);
            const day = d.getDay();
            const [, monthStr,] = d.toString().split(' ');
            return [`${monthStr} ${date}, ${year}`, acg.colors[day + 1]];
        }
        // 计算标签颜色
        function colorFromTarget(target) {
            let sum = 0;
            for (let char of target) {
                sum += char.codePointAt(0);
            }
            return acg.colors[sum % acg.colors.length];
        }
        let offset = 0;
        let total = 0;
        const articles = [];
        let zoneArticles = articles;
        function pullArticle() {
            const _pullArticle = pullArticle;
            pullArticle = () => { };
            function recover() { pullArticle = _pullArticle; }
            // 获取远程数据后回调
            {
                const data = [
                    {
                        "zone": "game",
                        "title": "某某文章标题某某文章标题",
                        "link": "#",
                        "date": "2017-04-03",
                        "image": "./static/image/acg_bgs/49647343_p0.png",
                        "abstracts": [
                            "大致介绍",
                            "继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话",
                            "继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话",
                        ],
                        "targets": [
                            "治愈",
                            "末世",
                            "纯爱",
                        ],
                    },
                    {
                        "zone": "music",
                        "title": "123某某文章标题某某文章标题",
                        "link": "#",
                        "date": "2017-04-03",
                        "image": "./static/image/acg_bgs/49647343_p0.png",
                        "abstracts": [
                            "大致介绍",
                            "继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话",
                            "继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话",
                        ],
                        "targets": [
                            "治愈",
                            "末世",
                            "纯爱",
                        ],
                    },
                    {
                        "zone": "game",
                        "title": "234某某文章标题某某文章标题",
                        "link": "#",
                        "date": "2017-04-03",
                        "image": "./static/image/acg_bgs/49647343_p0.png",
                        "abstracts": [
                            "大致介绍",
                            "继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话",
                            "继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话",
                        ],
                        "targets": [
                            "治愈",
                            "末世",
                            "纯爱",
                        ],
                    },
                    {
                        "zone": "music",
                        "title": "333某某文章标题某某文章标题",
                        "link": "#",
                        "date": "2017-04-03",
                        "image": "./static/image/acg_bgs/49647343_p0.png",
                        "abstracts": [
                            "大致介绍",
                            "继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话",
                            "继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话",
                        ],
                        "targets": [
                            "治愈",
                            "末世",
                            "纯爱",
                        ],
                    },
                    {
                        "zone": "game",
                        "title": "444某某文章标题某某文章标题",
                        "link": "#",
                        "date": "2017-04-03",
                        "image": "./static/image/acg_bgs/49647343_p0.png",
                        "abstracts": [
                            "大致介绍",
                            "继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话",
                            "继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话继续废话",
                        ],
                        "targets": [
                            "治愈",
                            "末世",
                            "纯爱",
                        ],
                    },
                ];
                data.forEach(({ zone, title, link, date, image, abstracts, targets, }) => {
                    const article = document.createElement('article');
                    const [dateStr, dayColor] = parseDate(date);
                    article.classList.add('article', dayColor);
                    article.innerHTML = `
                        <time>${dateStr}</time>
                                <h3>${title}</h3>
                                <section>
                                    <img src="${image}" alt="${title}" />
                                    <p>${abstracts.join('</p><p>')}</p>
                                </section>
                                <footer>
                                    ${targets.reduce((s, target) => s + `<i class="${colorFromTarget(target)}">${target}</i>`, '')}
                                </footer>
                    `;
                    article.addEventListener('click', e => location.href = link);
                    article.setAttribute('data-acg-zone', zone);
                    articles.push(article);
                });
                total = data.length;
                if (articles.length == total) {
                    recover = pullArticle = () => { };
                }
                recover();
            }
        }
        pullArticle();
        return {
            init(cards, left, right) {
                left.addEventListener("click", e => {
                    offset--;
                    if (offset < 0) {
                        offset = 0;
                    }
                    dataManager.load();
                });
                right.addEventListener("click", e => {
                    offset++;
                    if (offset > zoneArticles.length - 30) {
                        pullArticle();
                    } else if (offset >= total) {
                        offset = total - 1;
                    }
                    dataManager.load();
                });
                dataManager.cards = cards;
                dataManager.left = left;
                dataManager.right = right;
            },
            load(zone="") {
                dataManager.cards.forEach((card, index) => {
                    const article = zoneArticles[offset + index];
                    if (article == card.children[0]) {
                        return;
                    }
                    card.children[0] && card.removeChild(card.children[0]);
                    article && card.appendChild(article);
                });
            },
        }
    })();

    // 卡片模块
    {
        const cards = [];
        const centerDiv = document.querySelector('.-acg- .card>div.hover-center');
        if (innerWidth / innerHeight > 1) {
            // 横向，电脑
            cards.push((d => {
                d.classList.add('hover-left');
                centerDiv.before(d);
                return d;
            })(document.createElement('div')), (d => {
                return d;
            })(centerDiv), (d => {
                d.classList.add('hover-right')
                centerDiv.after(d);
                return d;
            })(document.createElement('div')));
        } else {
            // 竖向，手机
            cards.push(centerDiv);
        }
        cards.forEach((card, index) => {
            card.addEventListener('mouseover', e => {
                card.classList.remove("hover-left");
                card.classList.add("hover-center");
                card.classList.remove("hover-right");
                for (let left = 0; left < index; left++) {
                    cards[left].classList.add("hover-left");
                    cards[left].classList.remove("hover-center");
                    cards[left].classList.remove("hover-right");
                }
                for (let right = cards.length - 1; right > index; right--) {
                    cards[right].classList.remove("hover-left");
                    cards[right].classList.remove("hover-center");
                    cards[right].classList.add("hover-right");
                }
            });
        });
        // 左右滑动
        const left = document.querySelector('.-acg- .card>a:first-child');
        const right = document.querySelector('.-acg- .card>a:last-child');
        dataManager.init(cards, left, right);
        document.addEventListener("mousewheel", e => {
            let y = e.deltaY;
            let click = right.click.bind(right);
            if (e.deltaY < 0) {
                y = -y;
                click = left.click.bind(left);
            }
            for (y; y > 0; y -= 100) {
                click();
            }
        });
    }

    // 路由模块
    ((keys, fn) => keys.forEach(key => addEventListener(key, fn)))(
        ["load", "hashchange"],
        (() => {
            const style = ((style) => {
                document.head.appendChild(style);
                style.setAttribute("id", "acg-sticker-style");
                return style;
            })(document.createElement("style"));
            return (e) => {
                const hash = location.hash || "#";
                style.innerHTML = `
                    .-acg- .sticker>a[data-acg-hide-hash="${hash}"] {
                        display: none;
                    }
                    .-acg- .sticker>a[data-acg-show-hash="${hash}"] {
                        display: block;
                    }
                `;
                ({
                    "#": () => {
                        // 首页
                        dataManager.load();
                    },
                    "#anime": () => {
                        // 动画
                        dataManager.load("anime");
                    },
                    "#game": () => {
                        // 游戏
                        dataManager.load("game");
                    },
                    "#music": () => {
                        // 音乐
                        dataManager.load("music");
                    },
                    "#board": () => {
                        // 后台
                        console.log("b");
                    }
                })[hash]();
            }
        })()
    );
    // PATCH 修复UC浏览器onload提前触发问题
    document.readyState == 'complete' && window.dispatchEvent(new Event('load', { target: window }));

    // 初始方法
    {
        const acg = document.querySelector('.-acg-');
        const sticker = acg.querySelector('.sticker');
        const resize = sticker.querySelector('#resize:first-child');
        resize.addEventListener('click', e => sticker.classList.contains('fold') ? sticker.classList.remove('fold') : sticker.classList.add('fold'));
        // 存储折叠状态
        localStorage.getItem('isStickerFolded') && sticker.classList.add('fold');
        addEventListener("beforeunload", e => sticker.classList.contains('fold') ? localStorage.setItem('isStickerFolded', "true") : localStorage.removeItem('isStickerFolded'));
    }

    // Module defined.
    if (typeof module !== 'undefined' && module.exports) {
        module.exports = acg;
    } else if (typeof define === 'function' && define.amd) {
        define(function () { return acg; });
    } else {
        window.acg = acg;
    }
})();
