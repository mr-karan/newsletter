const BASE_URL = '/api';

const addNotification = (id) => {
    const alert = document.querySelector(id);
    alert.classList.remove("hide");
}

const addSubscription = async sub => {
    try {
        const res = await axios.post(`${BASE_URL}/create`, sub);
        const newSub = res.data;
        return newSub;
    } catch (e) {
        console.error(`Error creating new subscription ${e}`);
        return null;
    }
};

const form = document.querySelector('form');

const formEvent = form.addEventListener('submit', async event => {
    event.preventDefault();

    const btn = document.querySelector('#subscribe-btn');

    const email = document.querySelector('#email_addr').value;
    btn.disabled = true;
    btn.innerText = "Sending...";
    const sub = {
        email
    };

    const newSub = await addSubscription(sub);
    // if success
    if (newSub) {
        addNotification("#alert-success");
    } else {
        addNotification("#alert-error");
    }
    // reset state
    form.reset();
    btn.innerText = "Subscribe";
    btn.disabled = false;
});
