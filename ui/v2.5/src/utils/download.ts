const downloadFile = (url: string) => {
  const a = document.createElement("a");
  a.href = url;
  a.click();
};

export default downloadFile;
