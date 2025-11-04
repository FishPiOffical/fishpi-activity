(function() {
    'use strict';

    // 配置项
    const config = {
        apiUrl: 'https://fishpi-activity.aweoo.com/activity-api/recent', // API地址
        navSelector: '.nav-tabs', // 导航容器选择器
        linkHref: 'https://fishpi-activity.aweoo.com/', // 活动链接地址
        linkText: '活动', // 链接文本
        linkClass: '', // 链接class（可选）
        colors: {
            active: '#ff4757',      // 进行中 - 红色
            upcoming: '#ffa502'     // 即将开始 - 橙色
            // expired 不设置颜色，使用默认样式
        },
        showExpired: true // 是否显示已过期的活动标签
    };

    // 获取活动列表
    async function fetchActivities() {
        try {
            const response = await fetch(config.apiUrl);
            if (!response.ok) {
                throw new Error('Failed to fetch activities');
            }
            const data = await response.json();
            return data; // 直接返回 {active: [], upcoming: []}
        } catch (error) {
            console.error('Error fetching activities:', error);
            return { active: [], upcoming: [] };
        }
    }

    // 判断活动状态
    function getActivityStatus(data) {
        const activeActivities = data.active || [];
        const upcomingActivities = data.upcoming || [];

        if (activeActivities.length > 0) {
            return 'active';
        } else if (upcomingActivities.length > 0) {
            return 'upcoming';
        } else {
            return 'expired';
        }
    }

    // 插入导航标签
    function insertActivityNav(status) {
        const navTabs = document.querySelector(config.navSelector);
        if (!navTabs) {
            console.warn(`Element with selector "${config.navSelector}" not found`);
            return;
        }

        // 移除已存在的活动链接（防止重复插入）
        const existingLink = navTabs.querySelector('[data-fishpi-activity-nav]');
        if (existingLink) {
            existingLink.remove();
        }

        // 创建活动链接元素
        const activityLink = document.createElement('a');
        activityLink.href = config.linkHref;
        activityLink.textContent = config.linkText;
        activityLink.setAttribute('data-fishpi-activity-nav', 'true'); // 标记用于识别

        if (config.linkClass) {
            activityLink.className = config.linkClass;
        }

        // 根据状态设置样式
        if (status === 'active') {
            activityLink.style.color = config.colors.active;
            activityLink.style.fontWeight = 'bold';
            activityLink.setAttribute('data-activity-status', 'active');
        } else if (status === 'upcoming') {
            activityLink.style.color = config.colors.upcoming;
            activityLink.style.fontWeight = 'bold';
            activityLink.setAttribute('data-activity-status', 'upcoming');
        } else {
            activityLink.setAttribute('data-activity-status', 'expired');
        }
        // expired 状态不设置特殊样式，使用默认样式

        // 插入到第一个位置
        navTabs.insertBefore(activityLink, navTabs.firstChild);
    }

    // 初始化
    async function init() {
        const data = await fetchActivities();

        if (!data || (!data.active && !data.upcoming)) {
            return; // 没有活动，不显示导航标签
        }

        const status = getActivityStatus(data);

        // 根据配置决定是否显示过期活动标签
        if (status === 'expired' && !config.showExpired) {
            return;
        }

        insertActivityNav(status);
    }

    // 等待 DOM 加载完成后执行
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    // 暴露到全局，允许外部修改配置
    window.FishPiActivityNav = {
        config: config,
        init: init
    };
})();

