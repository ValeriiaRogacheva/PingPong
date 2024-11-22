<template>
  <div id="app">
    <h1>Ping Pong Game</h1>
    <div class="table" @mousemove="movePlayerPaddle">
      <div class="paddle player" :style="playerPaddleStyle"></div>
      <div class="paddle computer" :style="computerPaddleStyle"></div>
      <div class="ball" :style="ballStyle"></div>
    </div>
    <p>Player: {{ playerScore }} | Computer: {{ computerScore }}</p>
    <div class="controls">
      <button @click="changeSpeed('slow')">Slow</button>
      <button @click="changeSpeed('medium')">Medium</button>
      <button @click="changeSpeed('fast')">Fast</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      ballX: 250,
      ballY: 100,
      ballSpeedX: 3,
      ballSpeedY: 3,
      playerPaddleY: 80,
      computerPaddleY: 80,
      playerScore: 0,
      computerScore: 0,
      tableHeight: 200,
      paddleHeight: 40,
      paddleWidth: 10,
      ballSize: 20,
      computerTargetY: 80, // Целевая позиция ракетки компьютера
      computerSpeed: 1.2, // Скорость ракетки компьютера
      speedMultiplier: 1, // Множитель скорости
      computerErrorRate: 0.05, // Вероятность ошибки компьютера
    };
  },
  computed: {
    ballStyle() {
      return {
        left: `${this.ballX}px`,
        top: `${this.ballY}px`,
        position: "absolute",
        width: `${this.ballSize}px`,
        height: `${this.ballSize}px`,
        backgroundColor: "red",
        borderRadius: "50%",
        boxShadow: "0px 4px 8px rgba(0, 0, 0, 0.3)",
      };
    },
    playerPaddleStyle() {
      return {
        top: `${this.playerPaddleY}px`,
        left: "10px",
        position: "absolute",
        width: `${this.paddleWidth}px`,
        height: `${this.paddleHeight}px`,
        backgroundColor: "#007BFF",
        borderRadius: "5px",
        boxShadow: "0px 4px 8px rgba(0, 0, 0, 0.3)",
      };
    },
    computerPaddleStyle() {
      return {
        top: `${this.computerPaddleY}px`,
        right: "10px",
        position: "absolute",
        width: `${this.paddleWidth}px`,
        height: `${this.paddleHeight}px`,
        backgroundColor: "#28A745",
        borderRadius: "5px",
        boxShadow: "0px 4px 8px rgba(0, 0, 0, 0.3)",
      };
    },
  },
  methods: {
    movePlayerPaddle(event) {
      const tableTop = event.currentTarget.getBoundingClientRect().top;
      const mouseY = event.clientY - tableTop;
      this.playerPaddleY = Math.min(
          Math.max(mouseY - this.paddleHeight / 2, 0),
          this.tableHeight - this.paddleHeight
      );
    },
    startGame() {
      if (this.ballSpeedX === 0 && this.ballSpeedY === 0) {
        const randomDirection = Math.random() > 0.5 ? 1 : -1;
        this.ballSpeedX = randomDirection * 3 * this.speedMultiplier;
        this.ballSpeedY = randomDirection * 3 * this.speedMultiplier;
      }
    },
    updateGame() {
      // Обновляем позицию мяча
      this.ballX += this.ballSpeedX;
      this.ballY += this.ballSpeedY;

      // Проверяем столкновение с верхом и низом стола
      if (this.ballY <= 0 || this.ballY >= this.tableHeight - this.ballSize) {
        this.ballSpeedY *= -1;
      }

      // Проверяем столкновение с ракеткой игрока
      if (
          this.ballX <= this.paddleWidth + 10 &&
          this.ballY + this.ballSize >= this.playerPaddleY &&
          this.ballY <= this.playerPaddleY + this.paddleHeight
      ) {
        this.ballSpeedX = Math.abs(this.ballSpeedX);
      }

      // Проверяем столкновение с ракеткой компьютера
      if (
          this.ballX >= 500 - this.paddleWidth - 10 - this.ballSize &&
          this.ballY + this.ballSize >= this.computerPaddleY &&
          this.ballY <= this.computerPaddleY + this.paddleHeight
      ) {
        this.ballSpeedX = -Math.abs(this.ballSpeedX);
      }

      // Проверяем, кто пропустил мяч
      if (this.ballX < 0) {
        this.computerScore++;
        this.resetBall();
      } else if (this.ballX > 500) {
        this.playerScore++;
        this.resetBall();
      }

      // Движение ракетки компьютера
      this.moveComputerPaddle();
    },
    resetBall() {
      this.ballX = 250;
      this.ballY = 100;
      this.ballSpeedX = 0;
      this.ballSpeedY = 0;
    },
    moveComputerPaddle() {
      // Целевая позиция для компьютера
      const targetY = this.ballY - this.paddleHeight / 2;

      // Добавляем небольшую вероятность ошибки
      if (Math.random() < this.computerErrorRate) {
        this.computerTargetY = targetY + (Math.random() > 0.5 ? 15 : -15);
      } else {
        this.computerTargetY = targetY;
      }

      // Плавное движение к цели
      this.computerPaddleY +=
          (this.computerTargetY - this.computerPaddleY) * 0.2;

      // Ограничиваем движение ракетки в пределах стола
      this.computerPaddleY = Math.min(
          Math.max(this.computerPaddleY, 0),
          this.tableHeight - this.paddleHeight
      );
    },
    changeSpeed(level) {
      if (level === "slow") {
        this.speedMultiplier = 1;
        this.computerSpeed = 1.2;
        this.computerErrorRate = 0.05;
      }
      if (level === "medium") {
        this.speedMultiplier = 1.5;
        this.computerSpeed = 1.5;
        this.computerErrorRate = 0.03;
      }
      if (level === "fast") {
        this.speedMultiplier = 2;
        this.computerSpeed = 2;
        this.computerErrorRate = 0.01;
      }
    },
  },
  mounted() {
    window.addEventListener("click", this.startGame);
    setInterval(this.updateGame, 16);
  },
};
</script>

<style>
#app {
  text-align: center;
  font-family: Arial, sans-serif;
}

h1 {
  font-size: 24px;
  color: #333;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.1);
}

.table {
  width: 500px;
  height: 200px;
  border: 2px solid black;
  margin: 50px auto;
  position: relative;
  background: linear-gradient(135deg, #e8e8e8, #ffffff);
  border-radius: 10px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.2);
}

p {
  font-size: 18px;
  font-weight: bold;
}

button {
  margin: 5px;
  padding: 10px;
  font-size: 16px;
  cursor: pointer;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.2);
}

button:hover {
  background: #0056b3;
}
</style>