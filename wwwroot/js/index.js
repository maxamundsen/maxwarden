console.log("Previous - A powerful web codebase.");

// Helper functions

function toggleShowHide(element) {
    // Check current display style and toggle between 'none' and 'block'
    if (element.style.display === 'none' || element.style.display === "" ) {
        element.style.display = 'block';
    } else {
        element.style.display = 'none';
    }
}

function show(element) {
    element.style.display = 'block';
}

function hide(element) {
    element.style.display = 'none';
}

const onClickOutside = (element, callback) => {
    document.addEventListener('click', e => {
        if (!element.contains(e.target)) callback();
    });
};

const onClickOutsideOrEscape = (element, callback) => {
    document.addEventListener('click', e => {
        if (!element.contains(e.target)) callback();
    });

    document.addEventListener('keydown', (event) => {
        if (event.key === 'Escape') {
            const isNotCombinedKey = !(event.ctrlKey || event.altKey || event.shiftKey);
            if (isNotCombinedKey) {
                callback();
            }
        }
    });
};