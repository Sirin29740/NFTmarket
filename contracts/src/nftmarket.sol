// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol"; // 防止重入攻击（加分项）

contract NftMarket is ReentrancyGuard {
    struct Listing {
        address seller;
        address nftContract;
        uint256 tokenId;
        uint256 price;
        bool active;
    }

    // 记录所有挂单：listingId => Listing
    mapping(uint256 => Listing) public listings;
    uint256 private _listingIds;

    // 事件：方便后端和前端监听
    event ItemListed(uint256 indexed listingId, address indexed seller, uint256 price);
    event ItemSold(uint256 indexed listingId, address indexed buyer, uint256 price);

    // 1. 上架函数
    function listNFT(address nftContract, uint256 tokenId, uint256 price) external {
        require(price > 0, "Price must be > 0");

        // 检查：用户是否已经授权市场合约操作这个 NFT
        IERC721 nft = IERC721(nftContract);
        require(nft.getApproved(tokenId) == address(this) || nft.isApprovedForAll(msg.sender, address(this)), "Market not approved");

        _listingIds++;
        listings[_listingIds] = Listing(msg.sender, nftContract, tokenId, price, true);

        // 将 NFT 从卖家转入市场合约托管
        nft.transferFrom(msg.sender, address(this), tokenId);

        emit ItemListed(_listingIds, msg.sender, price);
    }

    // 2. 购买函数
    function buyNFT(uint256 listingId) external payable nonReentrant {
        Listing storage listing = listings[listingId];
        require(listing.active, "Item not for sale");
        require(msg.value >= listing.price, "Not enough ETH");

        listing.active = false; // 标记已售出

        // 转钱给卖家
        (bool success, ) = payable(listing.seller).call{value: listing.price}("");
        require(success, "Transfer failed");

        // 转 NFT 给买家
        IERC721(listing.nftContract).transferFrom(address(this), msg.sender, listing.tokenId);

        emit ItemSold(listingId, msg.sender, listing.price);
    }

    // 3. 下架函数
    function cancelListing(uint256 listingId) external {
        Listing storage listing = listings[listingId];
        require(msg.sender == listing.seller, "Not the seller");
        require(listing.active, "Not active");

        listing.active = false;
        IERC721(listing.nftContract).transferFrom(address(this), msg.sender, listing.tokenId);
    }
}