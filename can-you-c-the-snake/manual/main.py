from sklearn.datasets import load_iris
from sklearn.metrics import accuracy_score
from sklearn.model_selection import train_test_split
from sklearn.tree import DecisionTreeClassifier

# load dataset and split it into training and testing sets
X, y = load_iris(return_X_y=True)
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.3)

# define decision tree classifier model
clf = DecisionTreeClassifier()

# train the model
clf = clf.fit(X_train, y_train)

# model inference
pred = clf.predict(X_test)

# accuracy of the classifier is reported
acc = accuracy_score(y_test, pred)
print(f"The accuracy of the trained decision tree classifier is {acc*100:.1f} %")
