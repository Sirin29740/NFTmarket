// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";


contract MyNft is ERC721URIStorage,Ownable {
    uint256 private _tokenID;
    constructor() ERC721("MyNftMarket","mynft") Ownable(msg.sender) {}

    function mint(address to, string memory tokenURI) public onlyOwner returns(uint256){
        _tokenID += 1;
        uint256 newID = _tokenID;
        _safeMint(to, newID);
        _setTokenURI(newID, tokenURI);
        return newID;
    }

}
