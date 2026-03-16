import requests
import random

class StashInterface:
	port = ""
	url = ""
	headers = {
		"Accept-Encoding": "gzip, deflate, br",
		"Content-Type": "application/json",
		"Accept": "application/json",
		"Connection": "keep-alive",
		"DNT": "1"
		}

	def __init__(self, conn):
		self.port = conn['Port']
		scheme = conn['Scheme']

		self.url = scheme + "://localhost:" + str(self.port) + "/graphql"

		# Session cookie for authentication
		self.cookies = {}
		session_cookie = conn.get('SessionCookie')
		if isinstance(session_cookie, dict):
			session_value = session_cookie.get('Value')
			if session_value:
				self.cookies['session'] = session_value

	def __callGraphQL(self, query, variables = None):
		json = {}
		json['query'] = query
		if variables != None:
			json['variables'] = variables
		
		# handle cookies
		response = requests.post(self.url, json=json, headers=self.headers, cookies=self.cookies)
		
		if response.status_code == 200:
			result = response.json()
			# Properly handle GraphQL errors
			if result.get("errors"):
				raise Exception("GraphQL error: {}".format(result["errors"]))
			if result.get("data", None):
				return result.get("data")
			return None
		else:
			raise Exception("GraphQL query failed:{} - {}. Query: {}. Variables: {}".format(response.status_code, response.content, query, variables))

	def findTagIdWithName(self, name):
		query = """
query {
  allTags {
    id
    name
  }
}
		"""

		result = self.__callGraphQL(query)
		
		for tag in result["allTags"]:
			if tag["name"] == name:
				return tag["id"]
		return None

	def createTagWithName(self, name):
		query = """
mutation tagCreate($input:TagCreateInput!) {
  tagCreate(input: $input){
    id       
  }
}
"""
		variables = {'input': {
			'name': name
		}}

		result = self.__callGraphQL(query, variables)
		return result["tagCreate"]["id"]

	def destroyTag(self, id):
		query = """
mutation tagDestroy($input: TagDestroyInput!) {
  tagDestroy(input: $input)
}
"""
		variables = {'input': {
			'id': id
		}}

		self.__callGraphQL(query, variables)

	def findRandomSceneId(self):
		query = """
query findScenes($filter: FindFilterType!) {
  findScenes(filter: $filter) {
    count
    scenes {
      id
      tags {
        id
      }
    }
  }
}
"""
		
		variables = {'filter': {
			'per_page': 1,
			'sort': 'random'
		}}

		result = self.__callGraphQL(query, variables)

		if result["findScenes"]["count"] == 0:
			return None

		return result["findScenes"]["scenes"][0]

	def updateScene(self, sceneData):
		query = """
mutation sceneUpdate($input:SceneUpdateInput!) {
  sceneUpdate(input: $input) {
    id
  }
}
"""
		variables = {'input': sceneData}

		self.__callGraphQL(query, variables)

	# -------------------------------
	# Performer helpers for multiple images
	# -------------------------------

	def findPerformerIdWithName(self, name):
		query = """
query findPerformers($filter: FindFilterType!, $performer_filter: PerformerFilterType) {
  findPerformers(filter: $filter, performer_filter: $performer_filter) {
    count
    performers {
      id
      name
    }
  }
}
"""
		variables = {
			'filter': {'per_page': 1},
			'performer_filter': {
				'name': {'value': name, 'modifier': 'EQUALS'}
			}
		}
		result = self.__callGraphQL(query, variables)
		if not result or result["findPerformers"]["count"] == 0:
			return None
		return result["findPerformers"]["performers"][0]["id"]

	def getPerformer(self, performer_id):
		query = """
query getPerformer($id: ID!) {
  findPerformer(id: $id) {
    id
    name
    image_path
    images {
      id
      path
    }
  }
}
"""
		variables = {'id': performer_id}
		result = self.__callGraphQL(query, variables)
		return result.get("findPerformer") if result else None

	def setPerformerPrimaryImage(self, performer_id, image_id):
		# Updates the performer's primary/default image, if supported by the API
		query = """
mutation performerUpdate($input: PerformerUpdateInput!) {
  performerUpdate(input: $input) {
    id
  }
}
"""
		variables = {
			'input': {
				'id': performer_id,
				'primary_image_id': image_id
			}
		}
		self.__callGraphQL(query, variables)

	def addImageToPerformer(self, image_id, performer_id):
		# Associates an existing image with a performer
		query = """
mutation imageUpdate($input: ImageUpdateInput!) {
  imageUpdate(input: $input) {
    id
  }
}
"""
		variables = {
			'input': {
				'id': image_id,
				'performer_ids': [performer_id]
			}
		}
		self.__callGraphQL(query, variables)

	def removeImageFromPerformer(self, image_id, performer_id):
		# Disassociates an image from a performer
		# Fetch current performer_ids for image, remove one, and update
		# Note: Requires image performers query support.
		get_query = """
query getImage($id: ID!) {
  findImage(id: $id) {
    id
    performers {
      id
    }
  }
}
"""
		update_query = """
mutation imageUpdate($input: ImageUpdateInput!) {
  imageUpdate(input: $input) {
    id
  }
}
"""
		variables = {'id': image_id}
		result = self.__callGraphQL(get_query, variables)
		if not result or not result.get("findImage"):
			return
		current = result["findImage"].get("performers", [])
		current_ids = [p["id"] for p in current if p.get("id") and p["id"] != performer_id]
		update_vars = {'input': {'id': image_id, 'performer_ids': current_ids}}
		self.__callGraphQL(update_query, update_vars)

	def randomizePerformerPrimaryImage(self, performer_id):
		# Picks a random associated image and sets it as the primary/default image
		performer = self.getPerformer(performer_id)
		if not performer:
			return
		images = performer.get("images", []) or []
		if not images:
			return
		image = random.choice(images)
		self.setPerformerPrimaryImage(performer_id, image["id"])