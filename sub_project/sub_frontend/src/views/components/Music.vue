<template>
  <div class="music">
    <div class="left">
      <div class="item one"></div>
      <div class="item two"></div>
      <div class="item three"></div>
      <div class="item four"></div>
      <div class="item five"></div>
      <div class="item six"></div>
      <div class="item seven"></div>
      <div class="item eight"></div>
      <div class="item nine"></div>
      <div class="item ten"></div>
      <div class="item eleven"></div>
      <div class="item twelve"></div>
      <div class="item thirteen"></div>
      <div class="item fourteen"></div>
      <div class="item fifteen"></div>
      <div class="item sixteen"></div>
      <div class="item seventeen"></div>
      <div class="item eighteen"></div>
    </div>
    <div class="right">
      <div class="item eighteen"></div>
      <div class="item seventeen"></div>
      <div class="item sixteen"></div>
      <div class="item fifteen"></div>
      <div class="item fourteen"></div>
      <div class="item thirteen"></div>
      <div class="item twelve"></div>
      <div class="item eleven"></div>
      <div class="item ten"></div>
      <div class="item nine"></div>
      <div class="item eight"></div>
      <div class="item seven"></div>
      <div class="item six"></div>
      <div class="item five"></div>
      <div class="item four"></div>
      <div class="item three"></div>
      <div class="item two"></div>
      <div class="item one"></div>
    </div>
  </div>
</template>
<script lang="ts" setup></script>
<style scoped lang="scss">
.music {
  width: 405.7px;
  height: 74px;
  position: relative;
  .left {
    display: flex;
    align-items: center;
    width: 142.7px;
    position: absolute;
    left: 0;
    height: 74px;
    .item {
      width: 3.36px;
      margin-right: 5px;
      border-radius: 9px;
      background-color: rgba(122, 98, 245, 1);
    }
  }
  .right {
    display: flex;
    align-items: center;
    width: 142.7px;
    position: absolute;
    right: 0;
    height: 74px;
    .item {
      width: 3.36px;
      margin-right: 5px;
      border-radius: 9px;
      background-color: rgba(122, 98, 245, 1);
    }
  }
}
.one,
.four,
.seven,
.twelve,
.fifteen,
.eighteen {
  height: 26.09px;
  animation: music-animation 0.2s infinite alternate;
}
.two,
.five,
.eight,
.ten,
.thirteen,
.sixteen {
  height: 43.94px;
  animation: music-animation 0.5s infinite alternate;
}
.three,
.six,
.nine {
  height: 13.73px;
  animation: music-animation 0.7s infinite alternate;
}
.eleven,
.fourteen,
.seventeen {
  height: 74px;
  animation: music-animation 0.6s infinite alternate;
}

@keyframes music-animation {
  100% {
    height: 10px;
    border-radius: 9px;
    filter: contrast(2);
  }
}
</style>
