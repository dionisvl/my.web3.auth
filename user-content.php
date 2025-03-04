<?php

// start session
if (session_status() == PHP_SESSION_NONE) {
    session_start();
}

// example check web3 authorization
if(!empty($_SESSION['web3-wallet'])){
    echo 'Address User Wallet: '.$_SESSION['web3-wallet'];
    echo '<br>Show User Content';
} else {
    echo 'User Content Not Available';
}
echo '<br><a href="/web3-auth/web3-exit.php">Exit</a>';