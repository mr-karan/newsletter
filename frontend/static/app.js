const BASE_URL = 'https://jsonplaceholder.typicode.com';

export const addSubscription = async sub => {
    try {
        const res = await axios.post(`${BASE_URL}/subs`, sub);
        const newSub = res.data;

        console.log(`Added a new Subscription!`, newSub);

        return newSub;
    } catch (e) {
        console.error(`Error creating new subscription`, e);
    }
};

const form = document.querySelector('form');

const formEvent = form.addEventListener('submit', async event => {
    event.preventDefault();

    const email = document.querySelector('#email_addr').value;

    const sub = {
        email
    };

    const newSub = await addSubscription(sub);
    // addNotificationToDOM(newSub);
});
