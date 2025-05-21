function initializeEditor(editorQuery) {
    new EasyMDE({
        element: document.querySelector(editorQuery),
        forceSync: true
    })
}

function addMessage(content, timeout) {
    const messagePopupTmpl = document.querySelector('.message-popup-template')
    const messagePopup = messagePopupTmpl.cloneNode(true)
    messagePopup.classList.remove('message-popup-template')
    messagePopup.querySelector('.message-holder').innerHTML = "<span>" + content + "</span>"
    messagePopup.classList.remove('hide')
    messagePopupTmpl.parentNode.appendChild(messagePopup)
    if (timeout) {
        setTimeout(function () {
            messagePopup.classList.add('hide')
        }, timeout)
    }
}

function initializeForm(formQuery, submitUrl) {
    const form = document.querySelector(formQuery)

    function _prepareRequest() {
        const ret = {
            "original-page": form.querySelector('input[name="original-page"]').value,
            "page": form.querySelector('input[name="page"]').value,
            "content": form.querySelector('textarea[name="content"]').value,
            "comment": form.querySelector('input[name="comment"]').value,
        }
        if (window.turnstile) {
            ret.captcha = window.turnstile.getResponse()
        }
        return ret
    }

    form.querySelector('button[type="submit"]').addEventListener('click', async function (e) {
        e.preventDefault()
        const request = _prepareRequest()
        const response = await fetch(submitUrl, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(request),
        })

        const json = await response.json()
        addMessage(json.result.message)

        // Redirect if save was successful
        if (response.ok && json.result.is_success) {
            setTimeout(function () {
                window.location.href = "/wiki" + request.page + '.html'
            }, 1000)
        } else {
            if (window.turnstile) {
                window.turnstile.reset()
            }
        }
    })
}