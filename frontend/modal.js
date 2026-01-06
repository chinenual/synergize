// Adapted from https://jasonwatmore.com/post/2023/01/04/vanilla-js-css-modal-popup-dialog-tutorial-with-example

// open modal by id
export function openModal(id) {
    document.getElementById(id).classList.add('open');
    document.body.classList.add('modal-open');
}

// close currently open modal
export function closeModal() {
    document.querySelector('.modal.open').classList.remove('open');
    document.body.classList.remove('modal-open');
}

// redraw - maybe not needed post-bootstrap
export function updateModal(id) {
}

/**** 
window.addEventListener('load', function() {
    // close modals on background click
    document.addEventListener('click', event => {
        if (event.target.classList.contains('modal')) {
            closeModal();
        }
    });
});

*****/