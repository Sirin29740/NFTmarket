// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script} from "forge-std/Script.sol";
import {MyNft} from "../src/nft.sol"; // 确保路径指向你的合约

contract DeployMyNft is Script {
    function run() external {
        // 1. 从环境变量或配置文件读取私钥（或者在命令行通过 --private-key 传入）
        // uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        // 2. 开启广播：这之后的交易都会发送到链上
        vm.startBroadcast();

        // 3. 部署合约
        MyNft nft = new MyNft();

        // 4. 停止广播
        vm.stopBroadcast();
    }
}