class WebsocketClient {
  public websocket!: WebSocket;
  private wsUrl = 'ws://localhost:9111/ws';
  public isConnected = true;
  private messageQueue: string[] = [];

  private callbacks: {
    [eventName: string]: ((data?: any) => void)[];
  } = {};
  private callbacksConnections: ((isConnected: boolean) => void)[] = [];

  constructor() {
    this.initializeWebSocket();
  }

  public send(message: any) {
    if (!this.isConnected) {
      return this.messageQueue.push(message);
    }

    if (typeof message !== 'string') {
      try {
        message = JSON.stringify(message);
      } catch {}
    }
    this.websocket.send(message);
  }

  public onConnectionEvent(callback: (isConnected: boolean) => void) {
    this.callbacksConnections.push(callback);
  }
  public offConnectionEvent(callback: (isConnected: boolean) => void) {
    this.callbacksConnections = this.callbacksConnections.filter(
      c => c !== callback
    );
  }

  public on(eventName: string, callback: (data?: any) => void) {
    if (!this.callbacks[eventName]) {
      this.callbacks[eventName] = [];
    }
    this.callbacks[eventName].push(callback);
  }

  public off(eventName: string, callback: (data?: any) => void) {
    if (!this.callbacks[eventName]) {
      return;
    }
    this.callbacks[eventName] = this.callbacks[eventName].filter(
      c => c !== callback
    );
  }

  private initializeWebSocket() {
    if (!this.wsUrl) {
      return;
    }
    this.websocket = new WebSocket(this.wsUrl);
    this.websocket.onmessage = this.onMessage;
    this.websocket.onopen = this.onOpen;
    this.websocket.onclose = this.onClose;
    this.websocket.onerror = error => console.error(error);
  }

  private onMessage = (ev: MessageEvent) => {
    const { eventName, data } = JSON.parse(ev.data);

    if (this.callbacks[eventName]) {
      this.callbacks[eventName].forEach(cb => cb(data));
    }
  };

  private onOpen = () => {
    this.isConnected = true;

    this.callbacksConnections.forEach(cb => cb(true));

    this.messageQueue.forEach(message => this.send(message));
    this.messageQueue = [];
  };

  private onClose = () => {
    this.isConnected = false;

    this.callbacksConnections.forEach(cb => cb(false));

    setTimeout(() => {
      this.initializeWebSocket();
    }, 300);
  };
}

export const websocket = new WebsocketClient();
