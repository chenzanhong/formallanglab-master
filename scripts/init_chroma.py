#  pip install chromadb sentence-transformers
import json
import chromadb
from chromadb.utils import embedding_functions

# 加载问答对
with open("automata_qa.json", "r", encoding="utf-8") as f:
    qa_pairs = json.load(f)

# 初始化 Chroma 客户端（持久化到 ./chroma_db）
client = chromadb.PersistentClient(path="./chroma_db")

# 使用中文嵌入模型（需联网首次下载）
embedding_func = embedding_functions.SentenceTransformerEmbeddingFunction(
    model_name="BAAI/bge-large-zh-v1.5"
)

# 创建集合
collection = client.create_collection(
    name="automata_knowledge",
    embedding_function=embedding_func,
    metadata={"hnsw:space": "cosine"}
)

# 准备数据
documents = []
metadatas = []
ids = []

for i, qa in enumerate(qa_pairs):
    doc = f"问题：{qa['instruction']}\n答案：{qa['output']}"
    documents.append(doc)
    metadatas.append({"type": "qa", "topic": "automata"})
    ids.append(f"qa_{i}")

# 添加到集合
collection.add(
    documents=documents,
    metadatas=metadatas,
    ids=ids
)

print(f"成功导入 {len(documents)} 条知识到 ChromaDB！")

# python init_chroma.py