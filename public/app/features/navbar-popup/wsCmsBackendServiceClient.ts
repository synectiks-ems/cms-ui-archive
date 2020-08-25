const cfg = {
  URL: 'ws://100.81.5.26:4000/websocket/tracker/websocket',
};

const wsCmsBackendServiceSingletonClient = (() => {
  let instance: any;

  function createInstance() {
    const socket = new WebSocket(cfg.URL);
    return socket;
  }

  return {
    getInstance: () => {
      //   if (!instance) {
      instance = createInstance();
      //   }
      return instance;
    },
  };
})();

export default wsCmsBackendServiceSingletonClient;
