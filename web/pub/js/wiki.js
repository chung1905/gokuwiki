window.addEventListener('load', function () {
    requestAnimationFrame(function () {
        const codeBlocks = document.querySelectorAll('pre code');

        codeBlocks.forEach(function (codeBlock) {
            const container = codeBlock.parentNode.tagName === 'PRE' ?
                codeBlock.parentNode :
                codeBlock;

            container.style.position = 'relative';

            const copyButton = document.createElement('button');
            copyButton.innerText = 'Copy';
            copyButton.className = 'copy-button';
            copyButton.style.cssText = 'position: absolute; top: 5px; right: 5px; font-size: small; padding: 2px 5px; z-index: 10; opacity: 0.7;';

            copyButton.addEventListener('mouseover', function () {
                this.style.opacity = '1';
            });

            copyButton.addEventListener('mouseout', function () {
                this.style.opacity = '0.7';
            });

            copyButton.addEventListener('click', function () {
                const textToCopy = codeBlock.textContent;

                navigator.clipboard.writeText(textToCopy)
                    .then(() => {
                        const originalText = copyButton.innerText;
                        copyButton.innerText = 'Copied!';

                        setTimeout(() => {
                            copyButton.innerText = originalText;
                        }, 2000);
                    })
                    .catch(err => {
                        console.error('Failed to copy text: ', err);
                    });
            });

            container.appendChild(copyButton);
        });
    });
});