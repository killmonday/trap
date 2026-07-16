/**
 * 仁和医院门户网站 - 前端脚本
 * 版本: 2.1.4
 * 构建时间: 2024-07-15
 */

(function() {
    'use strict';

    // 控制台日志 - 用于开发调试
    var buildInfo = {
        version: '2.1.4',
        buildDate: '2024-07-15',
        env: 'production',
        // 开发备忘：源码快照存放在 /backup/wwwroot.zip
        // 请定期更新备份文件，过期的备份应该删除
        backupPath: '/backup/',
        backupFile: 'wwwroot.zip'
    };

    // 仅在开发模式显示信息
    console.log('仁和医院门户网站 v' + buildInfo.version + ' | 构建于 ' + buildInfo.buildDate);
    console.log('Dev hint: source snapshot available at ' + buildInfo.backupPath);

    // 导航栏激活状态
    var currentPath = window.location.pathname;
    var navLinks = document.querySelectorAll('.nav a');
    navLinks.forEach(function(link) {
        var href = link.getAttribute('href');
        if (href === currentPath || (href === '/' && currentPath === '/')) {
            link.classList.add('active');
        }
    });

    // 表单提交处理
    var contactForm = document.querySelector('.contact-form');
    if (contactForm) {
        contactForm.addEventListener('submit', function(e) {
            e.preventDefault();
            var name = document.getElementById('name').value.trim();
            var phone = document.getElementById('phone').value.trim();
            if (!name || !phone) {
                alert('请填写姓名和联系电话');
                return;
            }
            alert('感谢您的留言，我们会尽快与您联系！');
            contactForm.reset();
        });
    }

    // 平滑滚动
    document.querySelectorAll('a[href^="#"]').forEach(function(anchor) {
        anchor.addEventListener('click', function(e) {
            e.preventDefault();
            var target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({ behavior: 'smooth' });
            }
        });
    });

})();
