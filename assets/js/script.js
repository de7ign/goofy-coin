'esversion: 6'

/**
 *  createUser() takes the text from the respective input box
 *  and do a request to backend with the text for user creation
 */
function createUser() {
  // check: username can only be alphanumeric
  const alphaNumericPattern = /^([0-9]|[a-z]|[0-9a-z]+)$/i;
  const spacePattern = /\s+/g;

  const userName = document.getElementById("createUser").value;

  if (userName === "") {
    // set error message
    document.getElementById("createUserError").innerText =
      "please enter a username";
  } else if (spacePattern.test(userName)) {
    // set error message
    document.getElementById("createUserError").innerText = "no space allowed";
  } else if (!alphaNumericPattern.test(userName)) {
    // set error message
    document.getElementById("createUserError").innerText =
      "only numbers and alphabets are allowed";
  } else {
    // send request to backend
    document.getElementById("createUserError").innerText = "";
    let payload = {};
    payload.userName = userName;

    request("http://localhost:8080/user", payload)
      .then(response => {
        console.log(response);
      })
      .catch(err => {
        console.log(err);
      });
  }
}

/**
 *  createCoin() can only be called by 'goofy'.
 *
 *  createCoin() takes the text from the respective input box as the amount
 *  to create goofy coins and do a request to backend with the amount
 *  for creation of given goofy coins and transferred to goofy a/c
 */
function createCoin() {
  // check: current user is goofy or not
  const sel = document.getElementById("selectUser");
  if (sel.options[sel.selectedIndex].text !== "Goofy") {
    document.getElementById("createCoinError").innerText =
      "Only goofy can create coin";
    return;
  }

  // check: amount is numeric, non-empty and non-fractional
  const spacePattern = /\s+/g;
  const fractionPattern = /[.]+/g;
  const amount = document.getElementById("createCoin").value;
  if (amount === "") {
    // set error message
    document.getElementById("createCoinError").innerText = "enter an amount";
  } else if (spacePattern.test(amount)) {
    // set error message
    document.getElementById("createCoinError").innerText = "no space allowed";
  } else if (fractionPattern.test(amount)) {
    // set error message
    document.getElementById("createCoinError").innerText =
      "no fractions allowed";
  } else if (isNaN(amount)) {
    // set error message
    document.getElementById("createCoinError").innerText =
      "enter a numeric value";
  } else {
    // send request to backend
    document.getElementById("createCoinError").innerText = "";
  }
}

/**
 *  createTx() takes receiver name and amount to transfer
 *
 *  createTx() transfer goofy coins from current user in user section
 *  and transfer to the a/c of selected receiver user selected
 */
function createTx() {
  // check: receiver is selected or not
  const receiverSelection = document.getElementById("receiverPkeySelect");
  let receiver =
    receiverSelection.options[receiverSelection.selectedIndex].value;

  if (receiver === "0") {
    document.getElementById("payCoinError").innerText =
      "please select a receiver";
    return;
  }

  receiver = receiverSelection.options[receiverSelection.selectedIndex].text;
  const senderSelection = document.getElementById("selectUser");
  const sender = senderSelection.options[senderSelection.selectedIndex].text;

  // check: amount is numeric, non-empty and non-fractional
  const spacePattern = /\s+/g;
  const fractionPattern = /[.]+/g;
  const amount = document.getElementById("payAmount").value;
  if (amount === "") {
    // set error message
    document.getElementById("payCoinError").innerText = "enter an amount";
  } else if (spacePattern.test(amount)) {
    // set error message
    document.getElementById("payCoinError").innerText = "no space allowed";
  } else if (fractionPattern.test(amount)) {
    // set error message
    document.getElementById("payCoinError").innerText = "no fractions allowed";
  } else if (isNaN(amount)) {
    // set error message
    document.getElementById("payCoinError").innerText = "enter a numeric value";
  } else {
    // send request to backend
    document.getElementById("payCoinError").innerText = "";
  }
}

/**
 *  request() is a wrapper around axios post request call
 *
 *  @param {string} url   url where you want request to
 *  @param {Object} data  json data sent to url provided
 *
 *  @returns {pro} returns a promise for whichever callback is executed
 */
function request(url, data) {
  return axios
    .post(url, data)
    .then(function(response) {
      return response;
    })
    .catch(function(error) {
      return error.response;
    });
}
