function initializeEditor(editorQuery) {
    new EasyMDE({
        element: document.querySelector(editorQuery),
        forceSync: true
    })
}

function addMessage(content, timeout) {
    const messageHolder = document.querySelector(".message-holder")
    messageHolder.innerHTML = "<span>" + content + "</span>"
    messageHolder.classList.remove('hide')
}

function initializeForm(formQuery, submitUrl) {
    const form = document.querySelector(formQuery)

    function _prepareRequest() {
        const ret = {
            "page": form.querySelector('input[name="page"]').value,
            "content": form.querySelector('textarea[name="content"]').value,
            "comment": form.querySelector('input[name="comment"]').value,
        }
        if (turnstile) {
            ret.captcha = turnstile.getResponse()
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
        if (turnstile) {
            turnstile.reset()
        }
    })
}

function initializeTurnstile(renderContainer, turnstileSiteKey) {
}