// Adapted from
// https://jasonwatmore.com/post/2023/01/04/vanilla-js-css-modal-popup-dialog-tutorial-with-example
// changed to support the old bootstrap classes and attributes.

// open modal by id
export function openModal(id) {
  console.log('openModal()', id);
  let modal = document.getElementById(id);
  modal.classList.add('open');

  document.body.classList.add('modal-open');
}

// close currently open modal
export function closeModal() {
  console.trace('closeModal()');
  let modal = document.querySelector('.modal.open');
  if (modal) {
    modal.classList.remove('open');
  }
  document.body.classList.remove('modal-open');
}

// redraw - maybe not needed post-bootstrap.  TBD.
export function updateModal(id) {
  console.log('updateModal()', id);
}

let closeButtons = document.querySelectorAll('[data-bs-dismiss=modal]');
closeButtons.forEach((button) => {
        button.addEventListener('click', closeModal);
});

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